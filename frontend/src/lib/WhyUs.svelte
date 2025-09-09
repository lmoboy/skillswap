<script>
    import { Film, Headset, ShieldCheck } from "lucide-svelte";
</script>

<section
    class="w-full min-h-screen flex justify-center items-center bg-transparent gap-12 p-16 relative"
>
    <div class="flex gap-12 flex-wrap justify-center items-center w-full">
        <div
            class="relative backdrop-blur-lg bg-white/10 border border-indigo-400/30 rounded-2xl shadow-2xl shadow-indigo-900/30 hover:scale-105 transition-all duration-300 overflow-hidden p-8 w-72 h-80 flex flex-col justify-center items-center group"
        >
            <div class="flex items-center justify-center mb-4">
                <Film
                    class="w-12 h-12 stroke-indigo-400 group-hover:stroke-indigo-300 drop-shadow-[0_0_8px_rgba(99,102,241,0.7)] transition"
                />
            </div>
            <p
                class="mt-2 text-base font-semibold text-gray-200 text-center drop-shadow-sm"
            >
                Watch videos on the subject that interests you at any time made
                by people from elite universities.
            </p>
        </div>

        <div
            class="relative backdrop-blur-lg bg-white/10 border border-pink-400/30 rounded-2xl shadow-2xl shadow-pink-900/30 hover:scale-105 transition-all duration-300 overflow-hidden p-8 w-72 h-80 flex flex-col justify-center items-center group"
        >
            <div class="flex items-center justify-center mb-4">
                <ShieldCheck
                    class="w-12 h-12 stroke-pink-300 group-hover:stroke-pink-200 drop-shadow-[0_0_8px_rgba(236,72,153,0.7)] transition"
                />
            </div>
            <p
                class="mt-2 text-base font-semibold text-gray-200 text-center drop-shadow-sm"
            >
                Receive valid certifications upon exchanging skills with
                verified users*.
            </p>
        </div>

        <div
            class="relative backdrop-blur-lg bg-white/10 border border-cyan-400/30 rounded-2xl shadow-2xl shadow-cyan-900/30 hover:scale-105 transition-all duration-300 overflow-hidden p-8 w-72 h-80 flex flex-col justify-center items-center group"
        >
            <div class="flex items-center justify-center mb-4">
                <Headset
                    class="w-12 h-12 stroke-cyan-300 group-hover:stroke-cyan-200 drop-shadow-[0_0_8px_rgba(6,182,212,0.7)] transition"
                />
            </div>
            <p
                class="mt-2 text-base font-semibold text-gray-200 text-center drop-shadow-sm"
            >
                Watching videos can be boring, call the swapper instead! We
                support 1080p 60fps* video calls with sound!
            </p>
        </div>
    </div>
</section>
