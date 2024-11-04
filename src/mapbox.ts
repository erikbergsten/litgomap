import { LitElement, css, html } from 'lit'
import { customElement, property, query } from 'lit/decorators.js'
import { Map, Marker } from 'maplibre-gl';

@customElement("mp-map")
export class MpMax extends LitElement {

  @query('div')
  container

  firstUpdated() {
    console.log("i'm become rendered!", this.shadowRoot)
    window.map = new Map({
      container:  this.container,
      style: 'https://demotiles.maplibre.org/style.json', // stylesheet location
      center: [18, 59], // starting position [lng, lat]
      zoom: 6 // starting zoom
    })

    let marker = new Marker()
      .setLngLat([18.5, 59.8])
      .addTo(window.map);
  }

  render() {
    return html`
      <link class='mapstyle' rel='stylesheet' href='/node_modules/maplibre-gl/dist/maplibre-gl.css' />
      <div id='map-container'></div>
      <div id='rest'>
      </div>
    `
  }

  static styles = css`
    :host {
    }
    #map-container {
      width: 480px;
      height: 480px;
    }
    #rest {
      height: 240px;
      background-color: #fcc;
    }
  `
}
