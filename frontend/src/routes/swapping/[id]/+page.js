import { browser } from '$app/environment';

// @ts-ignore
export const load = async ({ params, fetch }) => {


    try {
        // @ts-ignore
        const u2id = params.id;
        // @ts-ignore
        const u1id = $auth.user.id;

        const res = await fetch(`/api/createChat?u1=${u1id}&u2=${u2id}`);
;
    } catch (error) {
        console.log(error)
    }
};