# RunCoach Pro - Assets

Colección de iconos SVG monocromáticos diseñados para RunCoach Pro.

## Paleta de Colores

- **Primary**: `#8b5cf6` (Violeta)
- **Secondary**: `#a78bfa` (Violeta claro)
- **Accent**: `#c084fc` (Violeta pastel)
- **Success**: `#10b981` (Verde)
- **Danger**: `#ef4444` (Rojo)
- **Warning**: `#f59e0b` (Ámbar)

## Gradiente Principal

```css
linear-gradient(135deg, #8b5cf6, #ec4899)
```

## Iconos Disponibles

### Logo Principal
- **logo.svg** - Logo completo con ritmo cardíaco y corredor

### Iconos de Métricas
- **icon-speed.svg** - Rayo (velocidad/ritmo)
- **icon-heart.svg** - Corazón (frecuencia cardíaca)
- **icon-heartbeat.svg** - Pulso cardíaco
- **icon-energy.svg** - Rayo relleno (energía/potencia)
- **icon-time.svg** - Reloj (tiempo/duración)
- **icon-chart.svg** - Gráfica (estadísticas)

### Iconos de Navegación
- **icon-arrow-right.svg** - Flecha derecha
- **icon-chevron-down.svg** - Chevron abajo
- **icon-location.svg** - Ubicación (rutas)
- **icon-calendar.svg** - Calendario (fechas)

### Iconos de Estado
- **icon-check.svg** - Check (completado)
- **icon-success.svg** - Éxito (guardado correctamente)
- **icon-info.svg** - Información
- **icon-refresh.svg** - Actualizar

### Iconos de Usuario
- **icon-user.svg** - Perfil de usuario

## Uso

```html
<!-- Inline en HTML -->
<img src="assets/icon-speed.svg" alt="Velocidad" width="24" height="24">

<!-- Como background en CSS -->
.icon {
  background-image: url('assets/icon-heart.svg');
  width: 24px;
  height: 24px;
}
```

## Características

- Todos los iconos son vectoriales (SVG)
- Tamaño base: 24x24px
- Stroke width: 2px
- Escalables sin pérdida de calidad
- Optimizados para tema dark
- Colores consistentes con la paleta de la app

## Personalización

Para cambiar el color de un icono, edita el atributo `stroke` o `fill` en el archivo SVG:

```svg
<path stroke="#8b5cf6" ... />  <!-- Color primario -->
<path stroke="#ef4444" ... />  <!-- Color de peligro -->
<path stroke="#10b981" ... />  <!-- Color de éxito -->
```
