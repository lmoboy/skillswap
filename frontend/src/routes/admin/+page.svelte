<script lang="ts">
   import { onMount } from 'svelte'
   import { auth } from '$lib/stores/auth'
   import { goto } from '$app/navigation'
   import LoadingSpinner from '$lib/components/common/LoadingSpinner.svelte'

   let activeTab = 'users'
   let users = []
   let skills = []
   let courses = []
   let loading = true
   let error = null

   let newSkill = { name: '', description: '' }

   import { browser } from '$app/environment'

   onMount(async () => {
      const unsubscribe = auth.subscribe((state) => {
         if (!state.loading) {
            if (!state.isAuthenticated || !state.user?.is_admin) {
               goto('/')
            } else {
               fetchData()
            }
         }
      })
      return unsubscribe
   })

   async function fetchData() {
      if (!browser) return
      loading = true
      try {
         if (activeTab === 'users') {
            const res = await fetch('/api/admin/users')
            users = await res.json()
         } else if (activeTab === 'skills') {
            const res = await fetch('/api/getSkills') // Public endpoint is fine for listing
            skills = await res.json()
         } else if (activeTab === 'courses') {
            const res = await fetch('/api/admin/courses')
            courses = await res.json()
         }
      } catch (e) {
         error = e.message
      } finally {
         loading = false
      }
   }

   $: if (activeTab) {
      fetchData()
   }

   async function deleteUser(id) {
      if (!confirm('Are you sure you want to delete this user?')) return
      const res = await fetch(`/api/admin/users/${id}`, { method: 'DELETE' })
      if (res.ok) fetchData()
   }

   async function deleteCourse(id) {
      if (!confirm('Are you sure you want to delete this course?')) return
      const res = await fetch(`/api/admin/courses/${id}`, { method: 'DELETE' })
      if (res.ok) fetchData()
   }

   async function deleteSkill(id) {
      if (!confirm('Are you sure you want to delete this skill?')) return
      const res = await fetch(`/api/admin/skills/${id}`, { method: 'DELETE' })
      if (res.ok) fetchData()
   }

   async function addSkill() {
      const res = await fetch('/api/admin/skills', {
         method: 'POST',
         headers: { 'Content-Type': 'application/json' },
         body: JSON.stringify(newSkill),
      })
      if (res.ok) {
         newSkill = { name: '', description: '' }
         fetchData()
      }
   }
</script>

<div class="admin-container">
   <h1>Admin Dashboard</h1>

   <div class="tabs">
      <button
         class:active={activeTab === 'users'}
         on:click={() => (activeTab = 'users')}>Users</button
      >
      <button
         class:active={activeTab === 'skills'}
         on:click={() => (activeTab = 'skills')}>Skills</button
      >
      <button
         class:active={activeTab === 'courses'}
         on:click={() => (activeTab = 'courses')}>Courses</button
      >
   </div>

   {#if loading}
      <LoadingSpinner />
   {:else if error}
      <p class="error">{error}</p>
   {:else if activeTab === 'users'}
      <table>
         <thead>
            <tr>
               <th>ID</th>
               <th>Username</th>
               <th>Email</th>
               <th>Admin</th>
               <th>Joined</th>
               <th>Actions</th>
            </tr>
         </thead>
         <tbody>
            {#each users as user}
               <tr>
                  <td>{user.id}</td>
                  <td>{user.username}</td>
                  <td>{user.email}</td>
                  <td>{user.is_admin ? 'Yes' : 'No'}</td>
                  <td>{new Date(user.created_at).toLocaleDateString()}</td>
                  <td>
                     <button
                        class="delete-btn"
                        on:click={() => deleteUser(user.id)}>Delete</button
                     >
                  </td>
               </tr>
            {/each}
         </tbody>
      </table>
   {:else if activeTab === 'skills'}
      <div class="add-skill">
         <h3>Add New Skill</h3>
         <input
            type="text"
            placeholder="Skill Name"
            bind:value={newSkill.name}
         />
         <input
            type="text"
            placeholder="Description"
            bind:value={newSkill.description}
         />
         <button on:click={addSkill}>Add Skill</button>
      </div>
      <table>
         <thead>
            <tr>
               <th>ID</th>
               <th>Name</th>
               <th>Description</th>
               <th>Actions</th>
            </tr>
         </thead>
         <tbody>
            {#each skills as skill}
               <tr>
                  <td>{skill.id}</td>
                  <td>{skill.name}</td>
                  <td>{skill.description}</td>
                  <td>
                     <button
                        class="delete-btn"
                        on:click={() => deleteSkill(skill.id)}>Delete</button
                     >
                  </td>
               </tr>
            {/each}
         </tbody>
      </table>
   {:else if activeTab === 'courses'}
      <table>
         <thead>
            <tr>
               <th>ID</th>
               <th>Title</th>
               <th>Instructor</th>
               <th>Skill</th>
               <th>Status</th>
               <th>Actions</th>
            </tr>
         </thead>
         <tbody>
            {#each courses as course}
               <tr>
                  <td>{course.id}</td>
                  <td>{course.title}</td>
                  <td>{course.instructor_name}</td>
                  <td>{course.skill_name}</td>
                  <td>{course.status}</td>
                  <td>
                     <button
                        class="delete-btn"
                        on:click={() => deleteCourse(course.id)}>Delete</button
                     >
                  </td>
               </tr>
            {/each}
         </tbody>
      </table>
   {/if}
</div>

<style>
   .admin-container {
      padding: 2rem;
      max-width: 1200px;
      margin: 0 auto;
      color: #1a202c; /* Dark text */
      background-color: #ffffff;
      min-height: 100vh;
   }

   h1 {
      font-size: 1.875rem;
      font-weight: 700;
      margin-bottom: 2rem;
      color: #111827;
      border-bottom: 2px solid #f3f4f6;
      padding-bottom: 1rem;
   }

   .tabs {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 2rem;
      background: #f9fafb;
      padding: 0.4rem;
      border-radius: 8px;
      width: fit-content;
   }

   .tabs button {
      padding: 0.5rem 1.25rem;
      border: none;
      background: transparent;
      cursor: pointer;
      border-radius: 6px;
      font-weight: 500;
      color: #4b5563;
      transition: all 0.2s;
   }

   .tabs button.active {
      background: #ffffff;
      color: #2563eb;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
   }

   .tabs button:hover:not(.active) {
      background: #f3f4f6;
      color: #111827;
   }

   table {
      width: 100%;
      border-collapse: separate;
      border-spacing: 0;
      margin-top: 1rem;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
      overflow: hidden;
   }

   th,
   td {
      padding: 1rem;
      text-align: left;
      border-bottom: 1px solid #e5e7eb;
   }

   th {
      background: #f8f9fa;
      font-weight: 600;
      color: #374151;
      text-transform: uppercase;
      font-size: 0.75rem;
      letter-spacing: 0.05em;
   }

   tr:last-child td {
      border-bottom: none;
   }

   tr:hover td {
      background-color: #f9fafb;
   }

   .delete-btn {
      background: #fee2e2;
      color: #dc2626;
      border: 1px solid #fecaca;
      padding: 0.4rem 0.8rem;
      border-radius: 6px;
      cursor: pointer;
      font-size: 0.875rem;
      font-weight: 500;
      transition: all 0.2s;
   }

   .delete-btn:hover {
      background: #dc2626;
      color: white;
   }

   .add-skill {
      margin-bottom: 2rem;
      padding: 1.5rem;
      background: #ffffff;
      border: 1px solid #e5e7eb;
      border-radius: 12px;
      display: flex;
      gap: 1rem;
      align-items: center;
      box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
   }

   .add-skill h3 {
      margin: 0;
      font-size: 1rem;
      margin-right: 1rem;
      color: #374151;
   }

   .add-skill input {
      padding: 0.6rem 1rem;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      flex: 1;
      font-size: 0.875rem;
   }

   .add-skill input:focus {
      outline: none;
      border-color: #2563eb;
      ring: 2px solid #bfdbfe;
   }

   .add-skill button {
      background: #2563eb;
      color: white;
      border: none;
      padding: 0.6rem 1.2rem;
      border-radius: 6px;
      cursor: pointer;
      font-weight: 500;
      transition: background 0.2s;
   }

   .add-skill button:hover {
      background: #1d4ed8;
   }

   .error {
      color: #dc2626;
      background: #fef2f2;
      padding: 1rem;
      border-radius: 8px;
      border: 1px solid #fee2e2;
   }
</style>
