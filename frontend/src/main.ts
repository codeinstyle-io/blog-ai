import './app.css';

import Posts from './apps/Posts.svelte';
import Pages from './apps/Pages.svelte';
import { Inity } from './lib/inity';

const Apps = {
  Posts,
  Pages,
};

Object.assign(window, {
  Inity,
  Apps,
});
