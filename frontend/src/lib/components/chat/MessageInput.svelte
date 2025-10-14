<script lang="ts">
    import { Send, Paperclip, Smile } from "lucide-svelte";

    type Props = {
        value?: string;
        placeholder?: string;
        disabled?: boolean;
        onSend?: (message: string) => void;
        onAttachment?: () => void;
        class?: string;
    };

    let {
        value = $bindable(""),
        placeholder = "Send a message...",
        disabled = false,
        onSend,
        onAttachment,
        class: className = ""
    }: Props = $props();

    function handleSubmit(event: Event) {
        event.preventDefault();

        if (!value.trim() || disabled) {
            return;
        }

        if (onSend) {
            onSend(value.trim());
        }

        value = "";
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            handleSubmit(event);
        }
    }

    function handleAttachmentClick() {
        if (onAttachment && !disabled) {
            onAttachment();
        }
    }
</script>

<div class="bg-white border-t border-gray-200 p-2 sm:p-4 {className}">
    <form onsubmit={handleSubmit} class="flex items-end gap-1.5 sm:gap-3">
        <button
            type="button"
            onclick={handleAttachmentClick}
            disabled={disabled}
            class="p-1.5 sm:p-2 rounded-full bg-gray-100 text-gray-600 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 hidden sm:flex"
            aria-label="Attach file"
        >
            <Paperclip class="h-4 w-4 sm:h-5 sm:w-5" />
        </button>

        <button
            type="button"
            disabled={disabled}
            class="p-1.5 sm:p-2 rounded-full bg-gray-100 text-gray-600 hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 hidden sm:flex"
            aria-label="Add emoji"
        >
            <Smile class="h-4 w-4 sm:h-5 sm:w-5" />
        </button>

        <div class="flex-1 relative">
            <textarea
                bind:value
                {placeholder}
                {disabled}
                onkeydown={handleKeydown}
                rows="1"
                maxlength="1000"
                class="w-full bg-gray-100 text-gray-900 p-2 sm:p-3 pr-12 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none disabled:opacity-50 disabled:cursor-not-allowed text-sm sm:text-base"
                style="max-height: 120px;"
            ></textarea>
            <span class="absolute bottom-1 right-2 text-xs text-gray-400">
                {value.length}/1000
            </span>
        </div>

        <button
            type="submit"
            disabled={disabled || !value.trim()}
            class="p-2 sm:p-3 rounded-lg bg-blue-500 text-white font-semibold hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 flex items-center gap-1.5 sm:gap-2"
        >
            <Send class="h-4 w-4 sm:h-5 sm:w-5" />
            <span class="hidden sm:inline">Send</span>
        </button>
    </form>
</div>
