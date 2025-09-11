// class AuthManager {
//     private name: string;
//     private email: string;
//     private picture: string;

//     constructor() {
//         this.email = ''
//         this.name = ''
//         this.picture = ''
//         this.getUserFromCookie()
//     }

//     public getUserInfo() {
//         return {
//             name: cookieStore.get("name") || this.name,
//             email: cookieStore.get("email") || this.email,
//             picture: cookieStore.get("picture") || this.picture,
//         };
//     }

//     public async getUserFromCookie() {
//         const response = await fetch("http://localhost:8080/api/cookieUser", {
//             credentials: "include",
//         }).then((response) => response.json()).catch((error) => console.log(error));
//         if (response == null || response == undefined) return;
//         const data = response;
//         // console.log(data);

//         if (data.error) {
//             console.log(data.error)
//             return
//         }
//         if (data.username == null || data.username == undefined ||
//             data.email == null || data.email == undefined
//         ) return;
//         this.name = data.username;
//         this.email = data.email;
//         this.picture = "bbcniggas";

//         cookieStore.set("name", this.name);
//         cookieStore.set("email", this.email);
//         cookieStore.set("picture", this.picture);
//     }
// }


// export default AuthManager;
