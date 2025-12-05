# Script para reemplazar alerts por showToast
$file = "d:\Metaphase07\Repositorios\TrainApp\trainapp\frontend\js\app.js"
$content = Get-Content $file -Raw -Encoding UTF8

# Reemplazos
$replacements = @{
    "alert\('✅ ¡Entreno guardado en el historial exitosamente!'\)" = "showToast('¡Entreno guardado en el historial exitosamente!', 'success')"
    "alert\('❌ Error al guardar el entreno'\)" = "showToast('Error al guardar el entreno', 'error')"
    "alert\('❌ Error de conexión'\)" = "showToast('Error de conexión', 'error')"
    "alert\('Error de conexión'\)" = "showToast('Error de conexión', 'error')"
    "alert\('Error al procesar la pregunta'\)" = "showToast('Error al procesar la pregunta', 'error')"
    "alert\('Error al generar el plan semanal'\)" = "showToast('Error al generar el plan semanal', 'error')"
    "alert\('Error al generar el plan de entrenamiento'\)" = "showToast('Error al generar el plan de entrenamiento', 'error')"
    "alert\('Error al generar el informe'\)" = "showToast('Error al generar el informe', 'error')"
    "alert\(`✅ Se importaron \$\{result.imported\} nuevos entrenamientos de Strava!`\)" = "showToast(`Se importaron `$`{result.imported`} nuevos entrenamientos de Strava!`, 'success')"
    "alert\('ℹ️ No hay entrenamientos nuevos para sincronizar'\)" = "showToast('No hay entrenamientos nuevos para sincronizar', 'info')"
    "alert\('❌ Error al sincronizar: ' \+ error\)" = "showToast('Error al sincronizar: ' + error, 'error')"
    "alert\('❌ Error de conexión al sincronizar con Strava'\)" = "showToast('Error de conexión al sincronizar con Strava', 'error')"
}

foreach ($pattern in $replacements.Keys) {
    $content = $content -replace $pattern, $replacements[$pattern]
}

Set-Content $file $content -Encoding UTF8 -NoNewline
Write-Host "✅ Reemplazos completados"
