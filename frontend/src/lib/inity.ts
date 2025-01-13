import { mount, type Component } from 'svelte';

type InityProps = {
  Component: Component;
  props: any | null;
};

export const Inity = {
  data: {} as {
    [key: string]: InityProps;
  },

  register(name: string, Component: Component, props: any | null): void {
    this.data[name] = {
      Component,
      props,
    };
  },

  attach(): void {
    for (const [name, entry] of Object.entries(this.data)) {
      const elements = document.querySelectorAll(`[x-inity="${name}"]`);
      const { Component } = entry as InityProps;

      elements.forEach((element) => {
        let props: any = {};
        let elementProps: string = element.getAttribute('x-props');

        if (elementProps) {
          try {
            props = JSON.parse(elementProps);
          } catch (e) {
            console.group('Inity');
            console.warn('Could not parse x-props attribute for element', element);
            console.warn('Props', elementProps);
            console.warn('Error', e);
            console.groupEnd();
          }
        }

        Object.assign(props, this.data[name].props);

        for (const [key, value] of Object.entries(this.data[name].props)) {
          if (value instanceof Function) {
            function wrapper(...args: any[]) {
              (value as Function)(...args, props);
            }
            props[key] = wrapper;
          }
        }

        mount(Component, { target: element, props });
      });
    }
  },
};
