<script lang="ts">
    let { user = "", message = "", isCurrentUser = false } = $props();
</script>

<div
    class="flex p-4 rounded-xl shadow-lg w-fit transition-all duration-300
         {isCurrentUser
        ? 'self-end bg-blue-600 text-white'
        : 'self-start bg-gray-700 text-gray-200'}"
>
    <div class="flex flex-col max-w-xs">
        {#if !isCurrentUser}
            <span class="text-sm font-semibold mb-1 opacity-80">
                {user}
            </span>
        {/if}

        <p class="text-base font-normal">
            {message}
        </p>
    </div>
</div>
