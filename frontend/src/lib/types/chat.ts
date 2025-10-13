/**
 * Chat and messaging type definitions
 */

export type User = {
    id: number | string;
    username: string;
    email?: string;
    profile_picture?: string;
};

export type Message = {
    id?: number;
    sender: User;
    content: string;
    timestamp: string;
    chat_id: number;
    read?: boolean;
};

export type ChatMeta = {
    id: number;
    user1_id: number;
    user2_id: number;
    created_at: string;
    user1_username: string;
    user1_profile_picture: string;
    user2_username: string;
    user2_profile_picture: string;
    last_message?: string;
    last_message_time?: string;
    unread_count?: number;
};

export type ChatWithMessages = ChatMeta & {
    messages: Message[];
};

export type WebSocketChatMessage = {
    type: 'post' | 'update' | 'new_message' | 'join' | 'leave';
    id?: number;
    chat_id?: number;
    user_id?: number | string;
    content?: string;
    message?: Message;
    [key: string]: any;
};

export type ChatRoom = {
    id: number;
    name: string;
    participants: User[];
    created_at: string;
};

export type TypingIndicator = {
    chat_id: number;
    user_id: number | string;
    username: string;
    is_typing: boolean;
};
