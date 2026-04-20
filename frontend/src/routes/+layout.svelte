<script>
   // @ts-nocheck

   // Import components
   import Header from '$lib/components/layout/Header.svelte'
   import Footer from '$lib/components/layout/Footer.svelte'

   import '../app.css'
   import { auth } from '$lib/stores/auth'
   import { checkAuth } from '$lib/api/auth'
   import Cookies from '$lib/components/ui/Cookies.svelte'

   // Get data from universal load function
   let { data } = $props()
   // Initialize auth store with data from load function
   // console.log($auth)
   // $effect(()=>{
   //   $auth
   //   console.log($auth)
   // })
   if(!$auth.isAuthenticated){
     checkAuth()
   }
   // if (data?.user) {
   //    auth.setUser({
   //       name: data.user.name,
   //       email: data.user.email,
   //       id: data.user.id,
   //       profile_picture: data.user.profile_picture || '',
   //    })
   // } else {
   //    auth.setLoading(false)
   // }
</script>

<div class="flex flex-col min-h-dvh relative">
   <Header class="shrink-0" />
   <main class="flex-1 w-full overflow-auto">
      <slot />
   </main>
   <Footer />
   <Cookies />
</div>
