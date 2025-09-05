<script>
    import { Film, Headset, ShieldCheck } from "lucide-svelte";
</script>

<section
    class="w-full min-h-screen flex justify-center items-center gap-12 p-16 bg-black relative"
>
    <div
        class="absolute bottom-[80%] left-[20%] bg-white/20 rounded-md rotate-[26deg] w-26 h-20 p-5 flex"
    ></div>
    <div
        class="absolute bottom-[30%] left-[50%] bg-white/20 rounded-md rotate-8 w-fit h-fit p-5 flex"
    ></div>
    <div
        class="absolute bottom-[60%] left-[70%] bg-white/20 rounded-md rotate-[56deg] w-20 h-8 flex"
    ></div>

    <div
        class="relative hover:scale-105 transition-all duration-300 overflow-hidden border border-gray-700 rounded-xl bg-white/10 p-8 w-64 h-64 flex flex-col justify-center items-center shadow-lg"
    >
        <Film class="w-10 h-10 stroke-white" />
        <p class="mt-4 text-base font-normal text-gray-300 text-center">
            Watch videos on the subject that interests you at any time made by
            people from elite universities.
        </p>
    </div>

    <div
        class="hover:scale-105 relative transition-all duration-300 overflow-hidden border border-gray-700 rounded-xl bg-white/10 p-8 w-64 h-64 flex flex-col justify-center items-center shadow-lg"
    >
        <!-- <GraduationCap  /> -->
        <ShieldCheck class="w-10 h-10 stroke-white" />
        <p class="mt-4 text-base font-normal text-gray-300 text-center">
            Receive valid certifications upon exchanging skills with verified
            users*.
        </p>
    </div>

    <div
        class="hover:scale-105 relative transition-all duration-300 overflow-hidden border border-gray-700 rounded-xl bg-white/10 p-8 w-64 h-64 flex flex-col justify-center items-center shadow-lg"
    >
        <Headset class="w-10 h-10 stroke-white" />
        <p class="mt-4 text-base font-normal text-gray-300 text-center">
            Watching videos can be boring, call the exchangee* instead! We support 1080p 60fps video calls with sound!
        </p>
    </div>
</section>
