import "../src/app.css";


import Posts from "../src/apps/Posts.svelte";
import Pages from "../src/apps/Pages.svelte";
import { Inity } from "../src/lib/inity";

Inity.register("posts", Posts, {
  title: "My First Post",
  slug: "my-first-post",
  excerpt: "This is my first post.",
  content: "Hello, world!",
  tags: ["tag1", "tag2", "tag3"],
  visible: true,
  savingState: "saved",
  onSubmit: (data, done, error) => {
    console.log("Submitted post data:", data);
    done("saved");
    error('Error');
  },
});

Inity.register("pages", Pages, {
  title: "My First Page",
  slug: "my-first-page",
  content: "Hello, world!",
  savingState: "saved",
  onSubmit: (data, done, error) => {
    console.log("Submitted page data:", data);
    done("saved");
    error('Error');
  },
});

Inity.attach();
