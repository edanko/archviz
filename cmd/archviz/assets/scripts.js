import 'htmx.org'

import Alpine from 'alpinejs'

// Add Alpine instance to window object.
window.Alpine = Alpine

// Start Alpine.
Alpine.start()

import Zoom from './zoom.js'

window.Zoom = Zoom

Zoom.init()

