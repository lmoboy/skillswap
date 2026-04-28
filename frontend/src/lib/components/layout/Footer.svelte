<script>


   import { Facebook, Linkedin, Mail, MapPin, Twitter } from 'lucide-svelte'
</script>

<footer class="bg-white border-t border-gray-200 mt-auto shrink-0">
   <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-10 md:py-12 flex flex-col sm:flex-row justify-between items-center gap-4">
      <p class="text-xs sm:text-sm text-gray-500 order-2 sm:order-1">
         {new Date().getFullYear()} SkillSwap&trade;
      </p>
      <div class="flex gap-6 order-1 sm:order-2">
         <a href="/privacy" class="text-xs sm:text-sm text-gray-500 hover:text-gray-900 transition-colors">Privacy Policy</a>
      </div>
   </div>
</footer>
