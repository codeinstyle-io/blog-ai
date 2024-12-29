import "./app.css";


import Posts from "./apps/Posts.svelte";
import { Inity } from "./lib/inity";
import type { SavingStates } from "./utils/types/common";


Inity.register('posts', Posts, { onSubmit: (element: Element, data: any, done: (savingState: SavingStates) => void) => {
  done('saving');
  fetch('http://localhost:8080/admin/api/posts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  }).then(
    (response) => response.json()
  ).then((data) => {
    console.log(data);
    done('saved');
  });
} })


document.addEventListener("DOMContentLoaded", () => Inity.attach());
