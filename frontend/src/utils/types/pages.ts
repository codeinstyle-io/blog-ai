
export  interface Pages {
    title?: string;
    slug?: string;
    content?: string;
    visible?: boolean;
    contentType?: string;
    onSubmit?: (data: any) => void;
  }