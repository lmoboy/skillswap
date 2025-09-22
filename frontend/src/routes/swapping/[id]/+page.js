import { browser } from '$app/environment';

export const load = async ({ params, fetch }) => {
    if (!browser) return {};
    
    try {
        // Get the current user's email from localStorage or session
        const userEmail = localStorage.getItem('userEmail');
        if (!userEmail) {
            throw new Error('User not authenticated');
        }

        // Get the other user's ID from the URL
        const otherUserId = params.id;

        // Get the other user's details
        const userResponse = await fetch(`/api/user?id=${otherUserId}`);
        if (!userResponse.ok) {
            throw new Error('Failed to fetch user details');
        }
        const otherUser = await userResponse.json();

        // Check if a chat already exists between these users
        const chatCheckResponse = await fetch('/api/checkChat', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user1Email: userEmail,
                user2Email: otherUser.email
            })
        });

        let chatId = null;
        if (chatCheckResponse.ok) {
            const chatData = await chatCheckResponse.json();
            chatId = chatData.chatId;
        }

        // If no chat exists, create a new one
        if (!chatId) {
            const createChatResponse = await fetch('/api/startChat', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: otherUser.email
                })
            });

            if (!createChatResponse.ok) {
                throw new Error('Failed to create chat');
            }

            const chatData = await createChatResponse.json();
            chatId = chatData.chatId;
        }

        // Get chat messages
        const messagesResponse = await fetch(`/api/messages?chatId=${chatId}`);
        const messages = messagesResponse.ok ? await messagesResponse.json() : [];

        return {
            otherUser: {
                id: otherUser.id,
                name: otherUser.username,
                email: otherUser.email,
                image: otherUser.profile_picture || `https://ui-avatars.com/api/?name=${encodeURIComponent(otherUser.username)}&background=random`
            },
            chatId,
            messages
        };
    } catch (error) {
        console.error('Error loading chat data:', error);
        // @ts-ignore
        return { error: error.message };
    }
};