<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { page } from '$app/stores';
  import { 
    Users, 
    BookOpen, 
    Target, 
    Activity, 
    LayoutDashboard, 
    Search, 
    Trash2, 
    ShieldAlert, 
    ShieldCheck, 
    Plus,
    ArrowLeft,
    RefreshCw,
    AlertTriangle,
    CheckCircle2,
    XCircle,
    ArrowUpRight,
    ArrowDownLeft,
    Coins
  } from 'lucide-svelte';

  // API functions
  async function fetchAdminStats() {
    const res = await fetch('/api/admin/stats', { credentials: 'include' });
    if (res.status === 401 || res.status === 403) {
      goto('/');
      return null;
    }
    return await res.json();
  }

  async function fetchUsers(search = '') {
    const url = search ? `/api/admin/users?search=${encodeURIComponent(search)}` : '/api/admin/users';
    const res = await fetch(url, { credentials: 'include' });
    if (res.status === 401 || res.status === 403) {
      goto('/');
      return null;
    }
    return await res.json();
  }

  async function toggleAdmin(userId: number, setAdmin: boolean) {
    const res = await fetch('/api/admin/user/toggle-admin', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, set_admin: setAdmin })
    });
    return await res.json();
  }

  async function deleteUser(userId: number) {
    const res = await fetch('/api/admin/user/delete', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId })
    });
    return await res.json();
  }

  async function fetchCourses() {
    const res = await fetch('/api/admin/courses', { credentials: 'include' });
    if (res.status === 401 || res.status === 403) {
      goto('/');
      return null;
    }
    return await res.json();
  }

  async function deleteCourse(courseId: number) {
    const res = await fetch('/api/admin/course/delete', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ course_id: courseId })
    });
    return await res.json();
  }

  async function fetchHealth() {
    const res = await fetch('/api/admin/health', { credentials: 'include' });
    return await res.json();
  }

  async function fetchSkills() {
    const res = await fetch('/api/admin/skills', { credentials: 'include' });
    if (res.status === 401 || res.status === 403) {
      goto('/');
      return null;
    }
    return await res.json();
  }

  async function addSkill(name: string, description: string) {
    const res = await fetch('/api/admin/skill/add', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, description })
    });
    return await res.json();
  }

  async function deleteSkill(skillId: number) {
    const res = await fetch('/api/admin/skill/delete', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ skill_id: skillId })
    });
    return await res.json();
  }

  // Swap API (Need to implement backend if not exist)
  async function updateUserSwaps(userId: number, amount: number) {
    const res = await fetch('/api/admin/user/swaps', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: userId, amount: amount })
    });
    return await res.json();
  }

  // State
  let stats: any = null;
  let users: any[] = [];
  let courses: any[] = [];
  let skills: any[] = [];
  let health: any = null;
  let loading = true;
  let activeTab = 'overview';
  let searchQuery = '';
  let message = '';
  let messageType: 'success' | 'error' = 'success';
  let itemToDelete: { type: 'user' | 'course' | 'skill', data: any } | null = null;
  let confirmDelete = false;
  let newSkillName = '';
  let newSkillDescription = '';
  
  // New Course Form State
  let showAddCourse = false;
  let newCourse = {
    title: '',
    description: '',
    difficulty_level: 'Beginner',
    skill_name: '',
    max_students: 10,
    status: 'Published'
  };

  onMount(async () => {
    if ($page.data.allowed === false) {
      goto('/auth/login');
      return;
    }

    const unsubscribe = auth.subscribe((state) => {
      if (!state.loading && !state.isAuthenticated) {
        goto('/auth/login');
      }
    });

    const user = await auth.waitForUser(5000).catch(() => null);
    if (!user) {
      goto('/auth/login');
      unsubscribe();
      return;
    }

    await loadData();
    unsubscribe();
  });

  async function loadData() {
    loading = true;
    const [statsData, usersData, coursesData, healthData, skillsData] = await Promise.all([
      fetchAdminStats(),
      fetchUsers(),
      fetchCourses(),
      fetchHealth(),
      fetchSkills()
    ]);

    stats = statsData;
    users = usersData?.users || [];
    courses = coursesData?.courses || [];
    skills = skillsData?.skills || [];
    health = healthData;
    loading = false;
  }

  function showMessage(msg: string, type: 'success' | 'error' = 'success') {
    message = msg;
    messageType = type;
    setTimeout(() => { message = ''; }, 3000);
  }

  async function handleToggleAdmin(user: any) {
    const result = await toggleAdmin(user.id, !user.is_admin);
    if (result.status === 'ok') {
      showMessage(`User ${user.username} admin status updated`);
      loadData();
    } else {
      showMessage(result.error || 'Failed to update user', 'error');
    }
  }

  async function handleSearch() {
    const data = await fetchUsers(searchQuery);
    users = data?.users || [];
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') {
      handleSearch();
    }
  }

  function confirmDeleteAction(type: 'user' | 'course' | 'skill', data: any) {
    itemToDelete = { type, data };
    confirmDelete = true;
  }

  async function executeDelete() {
    if (!itemToDelete) return;
    
    let result;
    if (itemToDelete.type === 'user') {
      result = await deleteUser(itemToDelete.data.id);
    } else if (itemToDelete.type === 'course') {
      result = await deleteCourse(itemToDelete.data.id);
    } else if (itemToDelete.type === 'skill') {
      result = await deleteSkill(itemToDelete.data.id);
    }

    if (result && result.status === 'ok') {
      showMessage(`${itemToDelete.type.charAt(0).toUpperCase() + itemToDelete.type.slice(1)} deleted`);
      loadData();
    } else {
      showMessage(result?.error || `Failed to delete ${itemToDelete.type}`, 'error');
    }
    
    confirmDelete = false;
    itemToDelete = null;
  }

  async function handleAddSkill() {
    if (!newSkillName.trim()) {
      showMessage('Skill name is required', 'error');
      return;
    }
    const result = await addSkill(newSkillName.trim(), newSkillDescription.trim());
    if (result.status === 'ok') {
      showMessage('Skill added');
      newSkillName = '';
      newSkillDescription = '';
      loadData();
    } else {
      showMessage(result.error || 'Failed to add skill', 'error');
    }
  }

  async function handleUpdateSwaps(userId: number, current: number, delta: number) {
    const result = await updateUserSwaps(userId, delta);
    if (result.status === 'ok') {
      showMessage(`User swaps updated`);
      loadData();
    } else {
      showMessage(result.error || 'Failed to update swaps', 'error');
    }
  }

  async function handleAddCourse() {
    if (!newCourse.title || !newCourse.skill_name) {
      showMessage('Title and Skill are required', 'error');
      return;
    }
    
    // Using multipart form as per backend AddCourse expectation
    const formData = new FormData();
    formData.append('title', newCourse.title);
    formData.append('description', newCourse.description);
    formData.append('difficulty_level', newCourse.difficulty_level);
    formData.append('skill_name', newCourse.skill_name);
    formData.append('max_students', newCourse.max_students.toString());
    formData.append('duration_minutes', '0'); // Required field in backend validation
    formData.append('modules', '[]'); // Empty modules for quick add
    
    const res = await fetch('/api/course/add', {
      method: 'POST',
      credentials: 'include',
      body: formData
    });
    
    const result = await res.json();
    if (res.ok) {
      showMessage('Course added successfully');
      showAddCourse = false;
      newCourse = { title: '', description: '', difficulty_level: 'Beginner', skill_name: '', max_students: 10, status: 'Published' };
      loadData();
    } else {
      showMessage(result.error || 'Failed to add course', 'error');
    }
  }
</script>

<svelte:head>
  <title>Admin Dashboard - SkillSwap</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 text-gray-900 font-sans w-full overflow-x-hidden">
  <!-- Header -->
  <header class="bg-white/80 backdrop-blur-md border-b border-gray-200 px-4 sm:px-8 py-4 sticky top-0 z-30 w-full">
    <div class="w-full flex items-center justify-between">

      <div class="flex items-center gap-3">
        <div class="p-2 bg-peach-100 rounded-lg">
          <ShieldAlert class="w-6 h-6 text-peach-600" />
        </div>
        <h1 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
          Admin Dashboard
        </h1>
      </div>
      <button
        onclick={() => goto('/')}
        class="flex items-center gap-2 px-4 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Back to Site</span>
      </button>
    </div>
  </header>

  <!-- Message Toast -->
  {#if message}
    <div class="fixed top-20 right-6 z-50 flex items-center gap-3 px-6 py-3 rounded-xl shadow-2xl animate-in slide-in-from-right fade-in duration-300 {messageType === 'success' ? 'bg-green-600 text-white' : 'bg-red-600 text-white'}">
      {#if messageType === 'success'}
        <CheckCircle2 class="w-5 h-5" />
      {:else}
        <XCircle class="w-5 h-5" />
      {/if}
      <span class="font-medium">{message}</span>
    </div>
  {/if}

  <!-- Confirm Delete Modal -->
  {#if confirmDelete}
    <div class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full border border-gray-100 overflow-hidden">
        <div class="p-6 text-center">
          <div class="w-16 h-16 bg-red-50 text-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
            <AlertTriangle class="w-8 h-8" />
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">Confirm Deletion</h3>
          <p class="text-gray-600">
            {#if itemToDelete?.type === 'user'}
              Are you sure you want to delete <strong>{itemToDelete.data.username}</strong>?
              <br/><span class="text-sm text-red-500 mt-2 block font-medium italic">This will permanently erase all their courses, chats, and messages.</span>
            {:else if itemToDelete?.type === 'course'}
              Are you sure you want to delete course <strong>{itemToDelete.data.title}</strong>?
            {:else if itemToDelete?.type === 'skill'}
              Are you sure you want to delete skill <strong>{itemToDelete.data.name}</strong>?
            {/if}
          </p>
        </div>
        <div class="flex border-t border-gray-100">
          <button
            onclick={() => { confirmDelete = false; itemToDelete = null; }}
            class="flex-1 px-6 py-4 text-gray-600 font-semibold hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onclick={executeDelete}
            class="flex-1 px-6 py-4 bg-red-600 text-white font-semibold hover:bg-red-700 transition-colors"
          >
            Yes, Delete
          </button>
        </div>
      </div>
    </div>
  {/if}

  <div class="w-full flex flex-col md:flex-row min-h-[calc(100vh-73px)]">
    <!-- Sidebar -->
    <aside class="w-full md:w-64 bg-white border-r border-gray-200 p-6 space-y-2 shrink-0">
      <nav class="space-y-1">
        <button
          onclick={() => activeTab = 'overview'}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all {activeTab === 'overview' ? 'bg-peach-50 text-peach-600 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
        >
          <LayoutDashboard class="w-5 h-5" />
          <span>Overview</span>
        </button>
        <button
          onclick={() => activeTab = 'users'}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all {activeTab === 'users' ? 'bg-peach-50 text-peach-600 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
        >
          <Users class="w-5 h-5" />
          <span>Users</span>
        </button>
        <button
          onclick={() => activeTab = 'courses'}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all {activeTab === 'courses' ? 'bg-peach-50 text-peach-600 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
        >
          <BookOpen class="w-5 h-5" />
          <span>Courses</span>
        </button>
        <button
          onclick={() => activeTab = 'skills'}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all {activeTab === 'skills' ? 'bg-peach-50 text-peach-600 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
        >
          <Target class="w-5 h-5" />
          <span>Skills</span>
        </button>
        <button
          onclick={() => activeTab = 'health'}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all {activeTab === 'health' ? 'bg-peach-50 text-peach-600 font-semibold' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'}"
        >
          <Activity class="w-5 h-5" />
          <span>System Health</span>
        </button>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-6 md:p-10 bg-gray-50/50 w-full">
      {#if loading}
        <div class="flex flex-col items-center justify-center h-64 gap-4">
          <RefreshCw class="w-10 h-10 text-peach-600 animate-spin" />
          <div class="text-gray-500 font-medium italic">Loading dashboard data...</div>
        </div>
      {:else}
        <!-- Overview Tab -->
        {#if activeTab === 'overview'}
          <div class="flex items-center justify-between mb-8">
            <h2 class="text-2xl font-bold text-gray-900">System Overview</h2>
            <button onclick={loadData} class="p-2 text-gray-500 hover:text-peach-600 hover:bg-white rounded-lg transition-all shadow-sm">
              <RefreshCw class="w-5 h-5" />
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
            <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 transition-transform hover:scale-[1.02]">
              <div class="flex items-center justify-between mb-4">
                <div class="p-2 bg-blue-50 text-blue-600 rounded-lg">
                  <Users class="w-5 h-5" />
                </div>
                <span class="text-xs font-bold text-green-600 bg-green-50 px-2 py-1 rounded-full">+{stats?.new_users_today || 0} today</span>
              </div>
              <div class="text-gray-500 text-sm font-medium mb-1">Total Users</div>
              <div class="text-3xl font-bold text-gray-900">{stats?.total_users || 0}</div>
            </div>

            <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 transition-transform hover:scale-[1.02]">
              <div class="flex items-center justify-between mb-4">
                <div class="p-2 bg-purple-50 text-purple-600 rounded-lg">
                  <ShieldCheck class="w-5 h-5" />
                </div>
              </div>
              <div class="text-gray-500 text-sm font-medium mb-1">Total Admins</div>
              <div class="text-3xl font-bold text-gray-900">{stats?.total_admins || 0}</div>
            </div>

            <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 transition-transform hover:scale-[1.02]">
              <div class="flex items-center justify-between mb-4">
                <div class="p-2 bg-orange-50 text-orange-600 rounded-lg">
                  <BookOpen class="w-5 h-5" />
                </div>
              </div>
              <div class="text-gray-500 text-sm font-medium mb-1">Total Courses</div>
              <div class="text-3xl font-bold text-gray-900">{stats?.total_courses || 0}</div>
            </div>

            <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 transition-transform hover:scale-[1.02]">
              <div class="flex items-center justify-between mb-4">
                <div class="p-2 bg-peach-50 text-peach-600 rounded-lg">
                  <Activity class="w-5 h-5" />
                </div>
              </div>
              <div class="text-gray-500 text-sm font-medium mb-1">Active Chats</div>
              <div class="text-3xl font-bold text-gray-900">{stats?.total_chats || 0}</div>
            </div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <div class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
              <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
                <Activity class="w-5 h-5 text-peach-600" />
                Platform Growth
              </h3>
              <div class="space-y-4">
                <div class="flex justify-between items-center p-3 hover:bg-gray-50 rounded-xl transition-colors">
                  <span class="text-gray-600 font-medium">New Users Today</span>
                  <span class="px-3 py-1 bg-green-50 text-green-700 rounded-full font-bold text-sm">+{stats?.new_users_today || 0}</span>
                </div>
                <div class="flex justify-between items-center p-3 hover:bg-gray-50 rounded-xl transition-colors">
                  <span class="text-gray-600 font-medium">New Users This Week</span>
                  <span class="px-3 py-1 bg-green-50 text-green-700 rounded-full font-bold text-sm">+{stats?.new_users_week || 0}</span>
                </div>
                <div class="flex justify-between items-center p-3 hover:bg-gray-50 rounded-xl transition-colors">
                  <span class="text-gray-600 font-medium">Registered Skills</span>
                  <span class="font-bold text-gray-900">{stats?.total_skills || 0}</span>
                </div>
                <div class="flex justify-between items-center p-3 hover:bg-gray-50 rounded-xl transition-colors">
                  <span class="text-gray-600 font-medium">Exchanged Messages</span>
                  <span class="font-bold text-gray-900">{stats?.total_messages || 0}</span>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
              <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
                <Target class="w-5 h-5 text-peach-600" />
                Quick Management
              </h3>
              <div class="grid grid-cols-1 gap-3">
                <button onclick={() => activeTab = 'users'} class="flex items-center justify-between p-4 bg-gray-50 hover:bg-peach-50 hover:text-peach-700 rounded-xl transition-all group font-medium">
                  <span>Manage User Accounts</span>
                  <Users class="w-5 h-5 text-gray-400 group-hover:text-peach-600" />
                </button>
                <button onclick={() => activeTab = 'courses'} class="flex items-center justify-between p-4 bg-gray-50 hover:bg-peach-50 hover:text-peach-700 rounded-xl transition-all group font-medium">
                  <span>Audit Course Catalog</span>
                  <BookOpen class="w-5 h-5 text-gray-400 group-hover:text-peach-600" />
                </button>
                <button onclick={() => activeTab = 'health'} class="flex items-center justify-between p-4 bg-gray-50 hover:bg-peach-50 hover:text-peach-700 rounded-xl transition-all group font-medium">
                  <span>System Health Check</span>
                  <Activity class="w-5 h-5 text-gray-400 group-hover:text-peach-600" />
                </button>
              </div>
            </div>
          </div>
        {/if}

        <!-- Users Tab -->
        {#if activeTab === 'users'}
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
            <h2 class="text-2xl font-bold text-gray-900">User Management</h2>
            
            <div class="flex gap-2 w-full md:w-auto">
              <div class="relative flex-1 md:w-80">
                <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <input
                  type="text"
                  bind:value={searchQuery}
                  onkeydown={handleKeydown}
                  placeholder="Search username or email..."
                  class="w-full pl-10 pr-4 py-2 bg-white border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all shadow-sm"
                />
              </div>
              <button
                onclick={handleSearch}
                class="px-6 py-2 bg-peach-600 text-white rounded-xl font-bold hover:bg-peach-700 transition-all shadow-sm shadow-peach-200"
              >
                Find
              </button>
              {#if searchQuery}
                <button
                  onclick={() => { searchQuery = ''; loadData(); }}
                  class="px-4 py-2 text-gray-500 hover:bg-white rounded-xl transition-all border border-transparent hover:border-gray-200"
                >
                  Clear
                </button>
              {/if}
            </div>
          </div>

          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left">
                <thead class="bg-gray-50 border-b border-gray-100">
                  <tr>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">User Profile</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Email Address</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Swaps</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Role</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                  {#each users as user}
                    <tr class="hover:bg-gray-50/50 transition-colors">
                      <td class="px-6 py-4">
                        <div class="flex items-center gap-3">
                          {#if user.profile_picture && user.profile_picture !== 'noPicture'}
                            <img src={user.profile_picture} alt="" class="w-10 h-10 rounded-full object-cover border-2 border-white shadow-sm" />
                          {:else}
                            <div class="w-10 h-10 rounded-full bg-peach-100 text-peach-700 flex items-center justify-center font-bold text-lg border-2 border-white shadow-sm">
                              {user.username?.[0]?.toUpperCase() || '?'}
                            </div>
                          {/if}
                          <div>
                            <div class="font-bold text-gray-900">{user.username}</div>
                            <div class="text-xs text-gray-400">ID: #{user.id} • Joined {user.created_at}</div>
                          </div>
                        </div>
                      </td>
                      <td class="px-6 py-4 text-gray-600 text-sm">{user.email}</td>
                      <td class="px-6 py-4 text-center">
                        <div class="flex items-center justify-center gap-2">
                           <button 
                             onclick={() => handleUpdateSwaps(user.id, user.swaps, -1)}
                             class="p-1 hover:bg-gray-200 rounded text-gray-500"
                             title="Decrease swaps"
                           >
                             <ArrowDownLeft class="w-3 h-3" />
                           </button>
                           <span class="font-bold text-gray-700 bg-gray-100 px-2 py-1 rounded-lg text-xs">{user.swaps}</span>
                           <button 
                             onclick={() => handleUpdateSwaps(user.id, user.swaps, 1)}
                             class="p-1 hover:bg-gray-200 rounded text-gray-500"
                             title="Increase swaps"
                           >
                             <ArrowUpRight class="w-3 h-3" />
                           </button>
                        </div>
                      </td>
                      <td class="px-6 py-4">
                        {#if user.is_admin}
                          <span class="inline-flex items-center gap-1.5 px-3 py-1 bg-red-50 text-red-600 rounded-full text-xs font-bold">
                            <ShieldCheck class="w-3 h-3" />
                            ADMIN
                          </span>
                        {:else}
                          <span class="inline-flex items-center gap-1.5 px-3 py-1 bg-gray-100 text-gray-600 rounded-full text-xs font-bold">
                            USER
                          </span>
                        {/if}
                      </td>
                      <td class="px-6 py-4 text-right">
                        <div class="flex justify-end gap-2">
                          <button
                            onclick={() => handleToggleAdmin(user)}
                            class="p-2 text-gray-400 hover:text-peach-600 hover:bg-peach-50 rounded-lg transition-all"
                            title={user.is_admin ? 'Remove Admin' : 'Make Admin'}
                          >
                            {#if user.is_admin}
                              <ShieldAlert class="w-5 h-5" />
                            {:else}
                              <ShieldCheck class="w-5 h-5" />
                            {/if}
                          </button>
                          <button
                            onclick={() => confirmDeleteAction('user', user)}
                            class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                            title="Delete User"
                          >
                            <Trash2 class="w-5 h-5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            {#if users.length === 0}
              <div class="p-20 text-center">
                <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Users class="w-8 h-8 text-gray-300" />
                </div>
                <div class="text-gray-400 font-medium">No users match your criteria</div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Courses Tab -->
        {#if activeTab === 'courses'}
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
            <h2 class="text-2xl font-bold text-gray-900">Course Management</h2>
            <button
              onclick={() => showAddCourse = !showAddCourse}
              class="flex items-center gap-2 px-6 py-3 bg-gray-900 text-white rounded-xl font-bold hover:bg-gray-800 transition-all shadow-lg"
            >
              <Plus class="w-5 h-5" />
              <span>{showAddCourse ? 'Cancel' : 'Add New Course'}</span>
            </button>
          </div>

          {#if showAddCourse}
            <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 mb-10 animate-in fade-in slide-in-from-top duration-300">
              <h3 class="text-lg font-bold text-gray-900 mb-6">Quick Course Creation</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Course Title</label>
                  <input
                    type="text"
                    bind:value={newCourse.title}
                    placeholder="e.g., Advanced Go Microservices"
                    class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                  />
                </div>
                <div>
                  <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Associated Skill</label>
                  <select
                    bind:value={newCourse.skill_name}
                    class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                  >
                    <option value="">Select a skill...</option>
                    {#each skills as skill}
                      <option value={skill.name}>{skill.name}</option>
                    {/each}
                  </select>
                </div>
                <div class="md:col-span-2">
                  <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Description</label>
                  <textarea
                    bind:value={newCourse.description}
                    placeholder="What will students learn in this course?"
                    rows="3"
                    class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                  ></textarea>
                </div>
                <div>
                  <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Difficulty Level</label>
                  <select
                    bind:value={newCourse.difficulty_level}
                    class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                  >
                    <option value="Beginner">Beginner</option>
                    <option value="Intermediate">Intermediate</option>
                    <option value="Advanced">Advanced</option>
                    <option value="Expert">Expert</option>
                  </select>
                </div>
                <div>
                  <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Max Students</label>
                  <input
                    type="number"
                    bind:value={newCourse.max_students}
                    min="1"
                    max="1000"
                    class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                  />
                </div>
              </div>
              <div class="mt-8 flex justify-end gap-3">
                <button
                  onclick={() => showAddCourse = false}
                  class="px-6 py-3 text-gray-500 font-bold hover:bg-gray-50 rounded-xl transition-all"
                >
                  Discard
                </button>
                <button
                  onclick={handleAddCourse}
                  class="px-10 py-3 bg-gray-900 text-white rounded-xl font-bold hover:bg-gray-800 transition-all shadow-lg"
                >
                  Create Course
                </button>
              </div>
            </div>
          {/if}

          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left">
                <thead class="bg-gray-50 border-b border-gray-100">
                  <tr>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Course Info</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Instructor</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Difficulty</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-center">Students</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Status</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                  {#each courses as course}
                    <tr class="hover:bg-gray-50/50 transition-colors">
                      <td class="px-6 py-4">
                        <div class="max-w-xs md:max-w-sm">
                          <div class="font-bold text-gray-900 truncate">{course.title}</div>
                          <div class="text-xs text-gray-400 truncate">{course.description || 'No description provided'}</div>
                        </div>
                      </td>
                      <td class="px-6 py-4">
                        <div class="flex items-center gap-2">
                          <div class="w-6 h-6 rounded-full bg-gray-200 flex items-center justify-center text-[10px] font-bold">
                            {course.instructor_name?.[0]?.toUpperCase() || '?'}
                          </div>
                          <span class="text-sm font-medium text-gray-700">{course.instructor_name || 'Unknown'}</span>
                        </div>
                      </td>
                      <td class="px-6 py-4 text-center">
                        <span class="px-2 py-1 bg-blue-50 text-blue-700 rounded-lg text-[10px] font-bold uppercase tracking-tight">
                          {course.difficulty_level}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-center font-bold text-gray-700 text-sm">
                        {course.current_students} <span class="text-gray-300 font-normal">/</span> {course.max_students}
                      </td>
                      <td class="px-6 py-4">
                        <span class="inline-flex items-center px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider
                          {course.status === 'Published' ? 'bg-green-100 text-green-700' :
                           course.status === 'Draft' ? 'bg-orange-100 text-orange-700' :
                           'bg-red-100 text-red-700'}">
                          {course.status}
                        </span>
                      </td>
                      <td class="px-6 py-4 text-right">
                        <button
                          onclick={() => confirmDeleteAction('course', course)}
                          class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                        >
                          <Trash2 class="w-5 h-5" />
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            {#if courses.length === 0}
              <div class="p-20 text-center">
                <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-4">
                  <BookOpen class="w-8 h-8 text-gray-300" />
                </div>
                <div class="text-gray-400 font-medium">No courses found in the system</div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Skills Tab -->
        {#if activeTab === 'skills'}
          <div class="flex items-center justify-between mb-8">
            <h2 class="text-2xl font-bold text-gray-900">Skills Library</h2>
          </div>

          <!-- Add Skill Form -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8 mb-10">
            <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
              <Plus class="w-5 h-5 text-peach-600" />
              Register New Skill
            </h3>
            <div class="flex flex-col lg:flex-row gap-4 items-end">
              <div class="flex-1 w-full">
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Skill Name</label>
                <input
                  type="text"
                  bind:value={newSkillName}
                  placeholder="e.g., Quantum Physics"
                  class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                />
              </div>
              <div class="flex-1 w-full">
                <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-2 ml-1">Short Description</label>
                <input
                  type="text"
                  bind:value={newSkillDescription}
                  placeholder="Briefly define this skill..."
                  class="w-full px-4 py-3 bg-gray-50 border border-gray-100 rounded-xl focus:outline-none focus:ring-2 focus:ring-peach-500/20 focus:border-peach-500 transition-all font-medium"
                />
              </div>
              <button
                onclick={handleAddSkill}
                class="w-full lg:w-auto px-8 py-3 bg-gray-900 text-white rounded-xl font-bold hover:bg-gray-800 transition-all shadow-lg"
              >
                Add Skill
              </button>
            </div>
          </div>

          <!-- Skills List -->
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left">
                <thead class="bg-gray-50 border-b border-gray-100">
                  <tr>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider w-16 text-center">ID</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Skill Name</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider">Description</th>
                    <th class="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-wider text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                  {#each skills as skill}
                    <tr class="hover:bg-gray-50/50 transition-colors">
                      <td class="px-6 py-4 text-center font-medium text-gray-400 text-sm">#{skill.id}</td>
                      <td class="px-6 py-4 font-bold text-gray-900">{skill.name}</td>
                      <td class="px-6 py-4 text-gray-500 text-sm">{skill.description || '-'}</td>
                      <td class="px-6 py-4 text-right">
                        <button
                          onclick={() => confirmDeleteAction('skill', skill)}
                          class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                        >
                          <Trash2 class="w-5 h-5" />
                        </button>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
            {#if skills.length === 0}
              <div class="p-20 text-center">
                <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Target class="w-8 h-8 text-gray-300" />
                </div>
                <div class="text-gray-400 font-medium">No skills registered yet</div>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Health Tab -->
        {#if activeTab === 'health'}
          <div class="flex items-center justify-between mb-8">
            <h2 class="text-2xl font-bold text-gray-900">System Health</h2>
            <button onclick={loadData} class="p-2 text-gray-500 hover:text-peach-600 hover:bg-white rounded-lg transition-all shadow-sm">
              <RefreshCw class="w-5 h-5" />
            </button>
          </div>

          <div class="grid gap-8">
            <div class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
              <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
                <Activity class="w-5 h-5 text-peach-600" />
                Database Engine
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="p-4 bg-gray-50 rounded-2xl border border-gray-100 flex items-center justify-between">
                  <span class="text-gray-600 font-medium">Status</span>
                  <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold
                    {health?.database === 'ok' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}">
                    <div class="w-2 h-2 rounded-full {health?.database === 'ok' ? 'bg-green-500 animate-pulse' : 'bg-red-500'}"></div>
                    {health?.database === 'ok' ? 'ONLINE' : 'CRITICAL ERROR'}
                  </span>
                </div>
                <div class="p-4 bg-gray-50 rounded-2xl border border-gray-100 flex items-center justify-between">
                  <span class="text-gray-600 font-medium">Last Sample</span>
                  <span class="text-gray-900 font-bold text-sm italic">{health?.timestamp ? new Date(health.timestamp).toLocaleTimeString() : 'N/A'}</span>
                </div>
              </div>
            </div>

            {#if health?.connections}
              <div class="bg-white rounded-2xl p-8 shadow-sm border border-gray-100">
                <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
                  <ShieldCheck class="w-5 h-5 text-peach-600" />
                  Connection Pool Performance
                </h3>
                <div class="grid grid-cols-2 lg:grid-cols-3 gap-6">
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Max Pool Size</div>
                    <div class="text-2xl font-bold text-gray-900">{health.connections.max_open}</div>
                  </div>
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Active Now</div>
                    <div class="text-2xl font-bold text-peach-600">{health.connections.open}</div>
                  </div>
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Work Load</div>
                    <div class="text-2xl font-bold text-gray-900">{health.connections.in_use}</div>
                  </div>
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Idle Slots</div>
                    <div class="text-2xl font-bold text-gray-900">{health.connections.idle}</div>
                  </div>
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Total Waits</div>
                    <div class="text-2xl font-bold text-gray-900">{health.connections.wait_count}</div>
                  </div>
                  <div class="p-4 border border-gray-100 rounded-2xl">
                    <div class="text-gray-400 text-xs font-bold uppercase tracking-wider mb-1">Latency Avg</div>
                    <div class="text-2xl font-bold text-gray-900">{health.connections.wait_duration}s</div>
                  </div>
                </div>
              </div>
            {/if}

            <div class="bg-gray-900 text-white rounded-2xl p-8 shadow-xl flex items-center justify-between overflow-hidden relative">
              <div class="relative z-10">
                <h4 class="text-xl font-bold mb-2">Need a full system reset?</h4>
                <p class="text-gray-400 text-sm max-w-md">Refreshing stats will pull the most recent data from the engine connection pool and recalculate growth metrics.</p>
              </div>
              <button
                onclick={loadData}
                class="relative z-10 px-8 py-4 bg-peach-600 text-white rounded-xl font-bold hover:bg-peach-700 transition-all shadow-lg flex items-center gap-2 group"
              >
                <RefreshCw class="w-5 h-5 group-hover:rotate-180 transition-transform duration-500" />
                Force Refresh
              </button>
              <!-- Decorative circles -->
              <div class="absolute -right-10 -bottom-10 w-40 h-40 bg-peach-500/20 rounded-full blur-3xl"></div>
              <div class="absolute -left-10 -top-10 w-32 h-32 bg-blue-500/10 rounded-full blur-3xl"></div>
            </div>
          </div>
        {/if}
      {/if}
    </main>
  </div>
</div>

<style>
  .bg-peach-50 { background-color: #fff7f5; }
  .bg-peach-100 { background-color: #ffedea; }
  .text-peach-600 { color: #f26444; }
  .text-peach-700 { color: #d44d2e; }
  .bg-peach-600 { background-color: #f26444; }
  .bg-peach-700 { background-color: #d44d2e; }
  .focus\:ring-peach-500\/20:focus { --tw-ring-color: rgba(242, 100, 68, 0.2); }
  .focus\:border-peach-500:focus { border-color: #f26444; }
  .shadow-peach-200 { --tw-shadow-color: rgba(242, 100, 68, 0.2); }
  .hover\:text-peach-600:hover { color: #f26444; }
  .hover\:text-peach-700:hover { color: #d44d2e; }
  .hover\:bg-peach-50:hover { background-color: #fff7f5; }
</style>
