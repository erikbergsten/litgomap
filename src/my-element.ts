import { LitElement, css, html } from 'lit'
import { customElement, property } from 'lit/decorators.js'
import litLogo from './assets/lit.svg'
import viteLogo from '/vite.svg'

@customElement('my-element')
export class MyElement extends LitElement {
  static shadowRootOptions = {...LitElement.shadowRootOptions, delegatesFocus: true , open: true }

  constructor() {
    super()
    const shadow = this.attachShadow({mode: 'open'})
    const div = document.createElement('div')
    div.innerHTML = `
      <div>
        <p> amy own stuff </p>
        <slot></slot>
      </div>
    `
    shadow.appendChild(div)
  }

  createRenderRoot() {
     return document.querySelector('#portal');
  }
  static styles = css`
    :host, #div {
      background-color: red;
      width: 100%;
      height: 100px;
      z-index: 99;
    }
  `
}
