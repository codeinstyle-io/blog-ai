import { mount, type Component } from "svelte";

type InityProps = {
    Component: Component,
    props: any | null
  }

export const Inity = {
    data: {} as {
      [key: string]: InityProps
    },

    register(name: string, Component: Component, props: any | null): void {
      this.data[name] = {
        Component,
        props,
      }
    },

    attach(): void {
      for (const [name, entry] of Object.entries(this.data)) {
        const elements = document.querySelectorAll(`[x-inity="${name}"]`)
        const { Component, props } = entry as InityProps

        elements.forEach((element) => {
          let props = element.getAttribute('x-props');

          if (props) {
            try {
              props = JSON.parse(props);
            } catch (e) {
              console.error(e);
            }
          }

          if (props) {
            Object.assign(props, this.data[name].props);
          }

          mount(Component, {target: element, props})
        })
      }
    }
  }