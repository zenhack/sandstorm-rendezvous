import RFB from '@novnc/novnc/core/rfb.js';


class VNCWindow extends HTMLElement {
  constructor() {
    super();
    this.setup();
  }

  setup() {
    const proto = window.location.protocol.replace('http', 'ws')
    const host = window.location.host;
    this.rfb = new RFB(this, proto + '//' + host + '/guest.socket');
    /*
    let src = undefined;
    if('src' in this.attributes) {
      src = this.attributes['src'].nodeValue;
      this.rfb = new RFB(this, src);
    }
    */
  }
}

customElements.define('vnc-window', VNCWindow);
