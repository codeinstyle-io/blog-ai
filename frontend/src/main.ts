import "./app.css";


import Posts from "./apps/Posts.svelte";
import { Inity } from "./lib/inity";

const Apps = {
  Posts,
};

Object.assign(window, {
  Inity,
  Apps,
});
