import "./app.css";


import Posts from "./apps/Posts.svelte";
import Pages from "./apps/Pages.svelte";
import { Inity } from "./lib/inity";

Inity.register('posts', Posts, { onSubmit: (data, done) => {
  console.log('Submitted post data:', data)
  done('saving');
  setTimeout(() => {
    done('saved');
  }, 1000);
} })

Inity.register('pages', Pages, { onSubmit: (data, done) => {
  console.log('Submitted page data:', data)
  done('saving');
  setTimeout(() => {
    done('saved');
  }, 1000);
} })

document.addEventListener("DOMContentLoaded", () => Inity.attach());
