
export  interface Posts {
    title?: string;
    slug?: string;
    tags?: string[];
    excerpt?: string;
    content?: string;
    visible?: boolean;
    publishDate?: Date;
    onSubmit?: (data: any) => void;
  }