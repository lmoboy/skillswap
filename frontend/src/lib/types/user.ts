export type User = {
    username: string;
    aboutme: string;
    email: string;
    location: string;
    profession: string;
    profile_picture: string;
    projects: Projects[];
    skills: Skills[];
    contacts: Contacts[];
};

export type Skills = {
    name: string;
    verified: boolean;
}

export type Contacts = {
    link: string;
    name: string;
    icon: string;
}

export type Projects = {
    link: string;
    name: string;
    description: string;
}