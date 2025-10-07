export type Course = {
    id: number;
    title: string;
    description: string;
    instructor_id: number;
    instructor_name: string;
    skill_id: number;
    skill_name: string;
    difficulty_level: 'Beginner' | 'Intermediate' | 'Advanced' | 'Expert';
    duration_hours: number;
    max_students: number;
    current_students: number;
    price: number;
    thumbnail_url: string;
    status: 'Draft' | 'Published' | 'Archived';
    created_at: string;
    updated_at: string;
    average_rating: number;
    review_count: number;
};

export type CourseModule = {
    id: number;
    course_id: number;
    title: string;
    description: string;
    order_index: number;
    created_at: string;
};

export type CourseEnrollment = {
    id: number;
    course_id: number;
    student_id: number;
    enrolled_at: string;
    completed_at: string;
    progress: number;
};

export type CourseReview = {
    id: number;
    course_id: number;
    student_id: number;
    student_name: string;
    rating: number;
    review_text: string;
    created_at: string;
};

export type CourseDetail = Course & {
    modules: CourseModule[];
    reviews: CourseReview[];
};
