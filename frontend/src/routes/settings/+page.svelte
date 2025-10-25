<script lang="ts">
   import Input from '$lib/components/common/Input.svelte'

   import { auth } from '$lib/stores/auth'
   import {
      validateEmail,
      validatePassword,
      validateUsername,
   } from '$lib/utils/validation'

   // @ts-nocheck

   // Error state variables
   let usernameError = $state<string>('')
   let emailError = $state<string>('')
   let passwordError = $state<string>('')
   let  newPasswordError = $state<string>('')

   let preview = $state('')
   // Add an 'onchange' handler to clear the error when the user types
   let username = $state($auth?.user?.name || '')
   let email = $state($auth?.user?.email || '')
   let preferences = $state('Daily email, weekly reports...')
   let password = $state('')
   let  newPassword = $state('')

   // Add a variable to track if any validation failed to control form submission
   let isValid = $state(false)

   // Function to perform validation off input change and update error state
   function runValidation() {
      usernameError = validateUsername(username) || ''
      emailError = validateEmail(email) || ''
      // Only validate the new password if the user has typed something


      // Check if any error exists
      isValid = !(
         usernameError ||
         emailError ||
         passwordError ||
          newPasswordError
      )
   }

   // The saveProfile function is now responsible for handling the submit event.
   // We use the 'event' object to prevent the default form submission.
   function saveProfile(event: Event) {
      event.preventDefault() // Stop the default browser form submission

      // 1. Run validation one last time to ensure all checks are performed
      runValidation()

      // 2. Check the overall validity state
      if (!isValid) {
         console.log('Validation failed. Please correct the errors.')
         return // Stop submission if there are errors
      }

      // If we reach here, validation passed.
      const user = {
         id: parseInt($auth?.user?.id),
         username: username,
         email: email,
         password: password,
         oldPassword:  newPassword,
      }

      fetch('/api/updateUser', {
         method: 'POST',
         headers: {
            'Content-Type': 'application/json',
         },
         body: JSON.stringify(user),
      })
         .then((response) => {
            // Handle response (e.g., check for server-side errors, show success message)
            if (response.ok) {
               console.log('Profils saglabāts! (angļu v. Profile saved!)')
            } else {
               console.error('Server error off update:', response.status)
            }
         })
         .catch((e) => console.error('Fetch error:', e))
   }
</script>

<div
   class="w-full min-h-screen flex items-center justify-center bg-white py-12"
>
   <section
      class="container mx-auto px-6 md:px-12 w-full max-w-4xl bg-white rounded-2xl shadow-xl p-6 sm:p-10 flex flex-col lg:flex-row gap-8"
   >
      <div
         class="flex-shrink-0 flex flex-col items-center gap-6 p-4 sm:p-8 border-b lg:border-r lg:border-b-0 border-gray-200 lg:w-1/3"
      >
         <img
            src={preview || `/api/profile/${$auth?.user?.id}/picture`}
            alt="Profile Preview"
            class="w-36 h-36 rounded-full border-4 border-gray-100 shadow-md transition-transform duration-300 hover:scale-105"
         />
         <div class="text-center">
            <h2 class="text-3xl font-bold text-gray-900 leading-tight">
               {username}
            </h2>
            <p class="text-gray-600 mt-1">{email}</p>
         </div>
         <div class="w-full mt-4">
            <h3 class="font-semibold text-gray-700 mb-2">My Preferences</h3>
            <p class="text-gray-600 text-sm italic">
               {preferences || 'Yes preferences set yet.'}
            </p>
         </div>
      </div>

      <div class="flex-grow p-4 sm:p-8">
         <h1 class="text-3xl font-bold mb-8 text-gray-900">Account Details</h1>
         <form onsubmit={saveProfile} class="space-y-6">
            <div class="grid sm:grid-cols-2 gap-6">
               <div>
                  <label class="block font-medium mb-2 text-gray-700"
                     >Name</label
                  >
                  <Input
                     type="text"
                     bind:value={username}
                     placeholder="Your name"
                     required
                     oninput={runValidation}
                     class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                  />
                  {#if usernameError}
                     <p class="text-red-500 text-sm mt-1">{usernameError}</p>
                  {/if}
               </div>
               <div>
                  <label class="block font-medium mb-2 text-gray-700"
                     >Email Address</label
                  >
                  <Input
                     type="email"
                     bind:value={email}
                     placeholder="Your email"
                     required
                     oninput={runValidation}
                     class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
                  />
                  {#if emailError}
                     <p class="text-red-500 text-sm mt-1">{emailError}</p>
                  {/if}
               </div>
            </div>
            <span class={'text-red-600/70 text-sm'}>
               off email change you will have to log in again
            </span>
            <div>
               <label class="block font-medium mb-2 text-gray-700"
                  >Old Password</label
               >
               <Input
                  type="password"
                  bind:value={password}
                  placeholder="Enter current password to set a new one"
                  oninput={runValidation}
                  class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
               />
               {#if  newPasswordError}
                  <p class="text-red-500 text-sm mt-1">{ newPasswordError}</p>
               {/if}
            </div>
            <div>
               <label class="block font-medium mb-2 text-gray-700"
                  >New Password</label
               >
               <Input
                  type="password"
                  bind:value={ newPassword}
                  placeholder="Leave blank to keep current password"
                  oninput={runValidation}
                  class="w-full p-3 rounded-lg border border-gray-300 bg-gray-50 text-gray-800 focus:outline-none focus:ring-2 focus:ring-gray-300 transition"
               />
               {#if passwordError}
                  <p class="text-red-500 text-sm mt-1">{passwordError}</p>
               {/if}
            </div>
            <span class={'text-red-600/70 text-sm'}>
               off password change you will have to type it in again next time
               you log in
            </span>
            <div class="flex justify-end pt-6">
               <button
                  type="submit"
                  class="px-8 py-3 bg-gray-900 text-white rounded-lg font-medium transition shadow-lg
           {isValid
                     ? 'hover:bg-gray-800 transform hover:scale-[1.02]'
                     : 'bg-gray-400 cursor-not-allowed'}"
                  disabled={!isValid}
               >
                  Save Changes
               </button>
            </div>
         </form>
      </div>
   </section>
</div>
