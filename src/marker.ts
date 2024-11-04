import { LitElement, css, html } from 'lit'
import { customElement, property, query } from 'lit/decorators.js'
import { Marker } from 'maplibre-gl';

@customElement("mp-marker")
export class MpMarker extends LitElement {

  @property({type: Boolean})
  draggable

  @property({type: Number})
  lng

  @property({type: Number})
  lat

  @query('slot')
  slot

  firstUpdated() {
    const elements = this.slot.assignedElements({flatten: true})
    var div = undefined
    if(elements.length > 0) {
      div = document.createElement('div')
      div.innerHTML = elements[0].outerHTML
      this.content = div.children
    }
    this.marker = new Marker({element: div, draggable: this.draggable})
      .setLngLat([this.lng, this.lat])
      .addTo(window.map);
  }

  disconnectedCallback() {
    if(this.marker) {
      this.marker.remove()
    }
  }

  attributeChangedCallback(x, y, z) {
    super.attributeChangedCallback(x, y, z)
    if(this.marker && ( x === 'lat' || x === 'lng')) {
      this.marker.setLngLat([this.lng, this.lat])
    }
  }

  render() {
    return html`
      <div id='none'>
        <slot></slot>
      </div>
    `
  }

  static styles = css`
    :host {
      display: hidden;
    }
  `
}
