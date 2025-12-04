// API Base URL - Automática según el entorno
const API_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8080/api'
    : `${window.location.origin}/api`;

// Estado global
let currentUser = null;
let allWorkouts = [];

// Inicializar la aplicación
document.addEventListener('DOMContentLoaded', () => {
    loadUser();
    loadWorkouts();
    setupEventListeners();
    checkStravaStatus();
    
    // Verificar si venimos del callback de Strava
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('strava') === 'connected') {
        alert('✅ ¡Strava conectado exitosamente! Sincronizando entrenamientos...');
        syncStrava();
        // Limpiar URL
        window.history.replaceState({}, document.title, window.location.pathname);
    }
});

// Setup Event Listeners
function setupEventListeners() {
    // Formulario de nuevo workout
    document.getElementById('workout-form').addEventListener('submit', handleWorkoutSubmit);
    
    // Formulario de plan de entrenamiento
    document.getElementById('plan-form').addEventListener('submit', handlePlanSubmit);
    
    // Formulario de informe
    document.getElementById('report-form').addEventListener('submit', handleReportSubmit);
    
    // Configurar fecha actual por defecto
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('date').value = now.toISOString().slice(0, 16);
}

// Switch entre métodos de entrada
function switchInputMethod(method) {
    // Actualizar botones
    document.querySelectorAll('.input-tab-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
    
    // Mostrar/ocultar formularios
    document.querySelectorAll('.input-method').forEach(el => el.classList.remove('active'));
    if (method === 'form') {
        document.getElementById('form-input').classList.add('active');
    } else {
        document.getElementById('image-input').classList.add('active');
    }
}

// Analizar workout desde imágenes
async function analyzeWorkoutFromImages() {
    const fileInput = document.getElementById('workout-images');
    const notes = document.getElementById('image-notes').value;
    
    if (!fileInput.files || fileInput.files.length === 0) {
        alert('Por favor, selecciona al menos una imagen');
        return;
    }
    
    try {
        showLoading(true);
        
        // Convertir imágenes a base64
        const imageURLs = [];
        for (let file of fileInput.files) {
            const base64 = await fileToBase64(file);
            imageURLs.push(base64);
        }
        
        const response = await fetch(`${API_URL}/workout-analysis-image`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                image_urls: imageURLs,
                notes: notes 
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            const resultDiv = document.getElementById('image-analysis-result');
            
            // Mostrar análisis con markdown y opción de guardar
            resultDiv.innerHTML = `
                <div class="analysis-section">
                    <h3>📊 Análisis del Entreno</h3>
                    <div class="markdown-content">${marked.parse(result.analysis)}</div>
                    
                    ${result.workout_data ? `
                        <div style="margin-top: 20px; padding: 15px; background: #e8f5e9; border-radius: 8px;">
                            <h4>✅ Datos extraídos del entreno</h4>
                            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 10px; margin-top: 10px;">
                                ${result.workout_data.type ? `<div><strong>Tipo:</strong> ${translateWorkoutType(result.workout_data.type)}</div>` : ''}
                                ${result.workout_data.distance ? `<div><strong>Distancia:</strong> ${result.workout_data.distance} km</div>` : ''}
                                ${result.workout_data.duration ? `<div><strong>Duración:</strong> ${result.workout_data.duration} min</div>` : ''}
                                ${result.workout_data.avg_pace ? `<div><strong>Ritmo:</strong> ${result.workout_data.avg_pace}</div>` : ''}
                                ${result.workout_data.avg_heart_rate ? `<div><strong>FC:</strong> ${result.workout_data.avg_heart_rate} bpm</div>` : ''}
                            </div>
                            <button onclick="saveExtractedWorkout(${JSON.stringify(result.workout_data).replace(/"/g, '&quot;')})" class="btn btn-primary" style="margin-top: 15px;">
                                💾 Guardar en el Historial
                            </button>
                        </div>
                    ` : ''}
                    
                    <!-- Chat de conversación -->
                    <div class="chat-container" style="margin-top: 20px;">
                        <div class="chat-messages" id="image-analysis-messages"></div>
                        <div class="chat-input-group">
                            <input type="text" id="image-analysis-question" placeholder="Pregunta algo sobre este análisis..." class="chat-input">
                            <button onclick="askAboutImageAnalysis()" class="btn btn-secondary">Enviar</button>
                        </div>
                    </div>
                </div>
            `;
            
            // Añadir el análisis como primer mensaje
            addChatMessage('assistant', result.analysis, 'image-analysis-messages');
            
            // Permitir enviar con Enter
            const input = document.getElementById('image-analysis-question');
            input.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    askAboutImageAnalysis();
                }
            });
            
            resultDiv.style.display = 'block';
            resultDiv.scrollIntoView({ behavior: 'smooth' });
        } else {
            alert('Error al analizar el entreno con imágenes');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Preguntar sobre el análisis de imagen
async function askAboutImageAnalysis() {
    const input = document.getElementById('image-analysis-question');
    const question = input.value.trim();
    
    if (!question) return;
    
    try {
        // Añadir pregunta del usuario al chat
        addChatMessage('user', question, 'image-analysis-messages');
        input.value = '';
        
        showLoading(true);
        
        // Enviar pregunta al backend
        const response = await fetch(`${API_URL}/workout-analysis-image`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                image_urls: [], // Ya no necesitamos las imágenes
                notes: '',
                question: question 
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Añadir respuesta del asistente al chat
            addChatMessage('assistant', result.analysis, 'image-analysis-messages');
            
            // Scroll al último mensaje
            const messagesDiv = document.getElementById('image-analysis-messages');
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        } else {
            alert('Error al procesar la pregunta');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Guardar workout extraído de la imagen
async function saveExtractedWorkout(workoutData) {
    try {
        showLoading(true);
        
        // Procesar fecha: intentar extraerla de workoutData o usar fecha actual
        let workoutDate = new Date().toISOString();
        
        if (workoutData.date) {
            // Si viene una fecha en formato ISO o timestamp
            try {
                workoutDate = new Date(workoutData.date).toISOString();
            } catch (e) {
                // Si falla, intentar parsear formatos comunes
                const dateMatch = workoutData.date.match(/(\d{4})-(\d{2})-(\d{2})/);
                if (dateMatch) {
                    workoutDate = new Date(`${dateMatch[0]}T12:00:00`).toISOString();
                }
            }
        }
        
        // Preparar el workout con valores por defecto para campos opcionales
        const workout = {
            user_id: currentUser ? currentUser.id : 1,
            date: workoutDate,
            type: workoutData.type || 'easy',
            distance: parseFloat(workoutData.distance) || 0,
            duration: parseInt(workoutData.duration) || 0,
            avg_pace: workoutData.avg_pace || '',
            avg_heart_rate: parseInt(workoutData.avg_heart_rate) || 0,
            avg_power: parseInt(workoutData.avg_power) || 0,
            cadence: parseInt(workoutData.cadence) || 0,
            elevation_gain: parseInt(workoutData.elevation_gain) || 0,
            calories: parseInt(workoutData.calories) || 0,
            feeling: workoutData.feeling || 'good',
            notes: workoutData.notes || 'Entreno importado desde captura del Apple Watch'
        };
        
        const response = await fetch(`${API_URL}/workouts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(workout)
        });
        
        if (response.ok) {
            alert('✅ ¡Entreno guardado en el historial exitosamente!');
            await loadWorkouts();
            
            // Cambiar a la pestaña de workouts para ver el entreno guardado
            showTab('workouts');
            document.querySelector('.tab-btn[onclick*="workouts"]').classList.add('active');
        } else {
            alert('❌ Error al guardar el entreno');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('❌ Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Convertir archivo a base64
function fileToBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

// Generar plan semanal contextual
async function generateWeeklyPlan() {
    try {
        showLoading(true);
        const response = await fetch(`${API_URL}/weekly-plan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Renderizar markdown
            const planContentDiv = document.getElementById('plan-content');
            planContentDiv.innerHTML = marked.parse(result.plan);
            
            // Mostrar el resultado y el chat
            const resultDiv = document.getElementById('plan-result');
            resultDiv.style.display = 'block';
            
            // Limpiar mensajes previos del chat
            document.getElementById('plan-messages').innerHTML = '';
            
            // Añadir el plan como primer mensaje del asistente
            addChatMessage('assistant', result.plan, 'plan-messages');
            
            // Scroll al resultado
            resultDiv.scrollIntoView({ behavior: 'smooth' });
        } else {
            alert('Error al generar el plan semanal');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Preguntar sobre el plan
async function askAboutPlan() {
    const input = document.getElementById('plan-question');
    const question = input.value.trim();
    
    if (!question) return;
    
    try {
        // Añadir pregunta del usuario al chat
        addChatMessage('user', question, 'plan-messages');
        input.value = '';
        
        showLoading(true);
        
        // Enviar pregunta al backend
        const response = await fetch(`${API_URL}/weekly-plan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ question: question })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Añadir respuesta del asistente al chat
            addChatMessage('assistant', result.plan, 'plan-messages');
            
            // Scroll al último mensaje
            const messagesDiv = document.getElementById('plan-messages');
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        } else {
            alert('Error al procesar la pregunta');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Añadir mensaje al chat
function addChatMessage(role, content, containerId) {
    const messagesDiv = document.getElementById(containerId);
    const messageDiv = document.createElement('div');
    messageDiv.className = `chat-message ${role}`;
    
    const label = document.createElement('strong');
    label.textContent = role === 'user' ? 'Tú:' : 'Entrenador:';
    
    const contentDiv = document.createElement('div');
    contentDiv.className = 'markdown-content';
    contentDiv.innerHTML = marked.parse(content);
    
    messageDiv.appendChild(label);
    messageDiv.appendChild(contentDiv);
    messagesDiv.appendChild(messageDiv);
}

// Permitir enviar con Enter
document.addEventListener('DOMContentLoaded', () => {
    const planQuestionInput = document.getElementById('plan-question');
    if (planQuestionInput) {
        planQuestionInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                askAboutPlan();
            }
        });
    }
});

// Tab Navigation
function showTab(tabName) {
    // Ocultar todos los tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Desactivar todos los botones
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Activar el tab seleccionado
    document.getElementById(tabName).classList.add('active');
    
    // Activar el botón correspondiente
    const targetBtn = document.querySelector(`.tab-btn[onclick*="${tabName}"]`);
    if (targetBtn) {
        targetBtn.classList.add('active');
    }
    
    // Recargar datos según el tab
    if (tabName === 'dashboard') {
        loadWorkouts();
    } else if (tabName === 'workouts') {
        displayAllWorkouts();
    }
}

// Cargar información del usuario
async function loadUser() {
    try {
        const response = await fetch(`${API_URL}/user`);
        const user = await response.json();
        currentUser = user;
        
        // Actualizar solo elementos que existen en el nuevo diseño
        const userName = document.getElementById('user-name');
        if (userName) {
            userName.textContent = user.name;
        }
    } catch (error) {
        console.error('Error cargando usuario:', error);
    }
}

// Cargar workouts
async function loadWorkouts() {
    try {
        const response = await fetch(`${API_URL}/workouts`);
        allWorkouts = await response.json() || [];
        
        updateDashboardStats();
        displayRecentWorkouts();
    } catch (error) {
        console.error('Error cargando workouts:', error);
        allWorkouts = [];
    }
}

// Actualizar estadísticas del dashboard
function updateDashboardStats() {
    const totalWorkouts = allWorkouts.length;
    const totalDistance = allWorkouts.reduce((sum, w) => sum + (w.distance || 0), 0);
    const totalTime = allWorkouts.reduce((sum, w) => sum + (w.duration || 0), 0);
    
    // Actualizar estadísticas si existen los elementos
    const totalWorkoutsEl = document.getElementById('total-workouts');
    const totalDistanceEl = document.getElementById('total-distance');
    const totalTimeEl = document.getElementById('total-time');
    
    if (totalWorkoutsEl) totalWorkoutsEl.textContent = totalWorkouts;
    if (totalDistanceEl) totalDistanceEl.textContent = `${totalDistance.toFixed(1)} km`;
    if (totalTimeEl) totalTimeEl.textContent = `${totalTime} min`;
    
    // Actualizar header stats
    const headerStats = document.getElementById('header-stats');
    if (headerStats) {
        if (totalWorkouts === 0) {
            headerStats.textContent = '¡Comienza tu primer entreno!';
        } else {
            const thisWeek = getThisWeekWorkouts();
            const weekDistance = thisWeek.reduce((sum, w) => sum + (w.distance || 0), 0);
            headerStats.textContent = `${thisWeek.length} entrenos esta semana • ${weekDistance.toFixed(1)} km`;
        }
    }
}

// Obtener workouts de esta semana
function getThisWeekWorkouts() {
    const now = new Date();
    const startOfWeek = new Date(now);
    startOfWeek.setDate(now.getDate() - now.getDay() + 1); // Lunes
    startOfWeek.setHours(0, 0, 0, 0);
    
    return allWorkouts.filter(w => {
        const workoutDate = new Date(w.date);
        return workoutDate >= startOfWeek;
    });
}

// Mostrar últimos workouts
function displayRecentWorkouts() {
    const container = document.getElementById('recent-workouts-list');
    const recent = allWorkouts.slice(0, 5);
    
    if (recent.length === 0) {
        container.innerHTML = '<p>No hay entrenamientos registrados aún.</p>';
        return;
    }
    
    container.innerHTML = recent.map(workout => createWorkoutCard(workout)).join('');
}

// Mostrar todos los workouts
function displayAllWorkouts() {
    const container = document.getElementById('workouts-list');
    
    if (allWorkouts.length === 0) {
        container.innerHTML = '<p>No hay entrenamientos registrados aún.</p>';
        return;
    }
    
    container.innerHTML = allWorkouts.map(workout => createWorkoutCard(workout, true)).join('');
}

// Filtrar workouts por tipo y mes
function filterWorkouts() {
    const typeFilter = document.getElementById('filter-type').value;
    const monthFilter = parseInt(document.getElementById('filter-month').value);
    
    let filtered = allWorkouts;
    
    // Filtrar por tipo
    if (typeFilter !== 'all') {
        filtered = filtered.filter(w => w.type === typeFilter);
    }
    
    // Filtrar por mes
    if (monthFilter !== 'all') {
        const now = new Date();
        const cutoffDate = new Date(now);
        cutoffDate.setMonth(cutoffDate.getMonth() - monthFilter);
        
        filtered = filtered.filter(w => {
            const workoutDate = new Date(w.date);
            return workoutDate >= cutoffDate;
        });
    }
    
    // Mostrar workouts filtrados
    const container = document.getElementById('workouts-list');
    
    if (filtered.length === 0) {
        container.innerHTML = '<p>No se encontraron entrenamientos con estos filtros.</p>';
        return;
    }
    
    container.innerHTML = filtered.map(workout => createWorkoutCard(workout, true)).join('');
}

// Crear tarjeta de workout
function createWorkoutCard(workout, showAnalyzeBtn = false) {
    const date = new Date(workout.date);
    const formattedDate = date.toLocaleDateString('es-ES', { 
        year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' 
    });
    
    return `
        <div class="workout-item">
            <div class="workout-header">
                <span class="workout-date">${formattedDate}</span>
                <span class="workout-type">${translateWorkoutType(workout.type)}</span>
            </div>
            <div class="workout-stats">
                <div class="workout-stat">
                    <span class="workout-stat-label">Distancia</span>
                    <span class="workout-stat-value">${workout.distance} km</span>
                </div>
                <div class="workout-stat">
                    <span class="workout-stat-label">Duración</span>
                    <span class="workout-stat-value">${workout.duration} min</span>
                </div>
                <div class="workout-stat">
                    <span class="workout-stat-label">Ritmo</span>
                    <span class="workout-stat-value">${workout.avg_pace}</span>
                </div>
                <div class="workout-stat">
                    <span class="workout-stat-label">FC Media</span>
                    <span class="workout-stat-value">${workout.avg_heart_rate || '-'} bpm</span>
                </div>
                ${workout.avg_power ? `
                <div class="workout-stat">
                    <span class="workout-stat-label">Potencia</span>
                    <span class="workout-stat-value">${workout.avg_power} W</span>
                </div>` : ''}
                ${workout.cadence ? `
                <div class="workout-stat">
                    <span class="workout-stat-label">Cadencia</span>
                    <span class="workout-stat-value">${workout.cadence} ppm</span>
                </div>` : ''}
                ${workout.elevation_gain ? `
                <div class="workout-stat">
                    <span class="workout-stat-label">Desnivel +</span>
                    <span class="workout-stat-value">${workout.elevation_gain} m</span>
                </div>` : ''}
                <div class="workout-stat">
                    <span class="workout-stat-label">Sensación</span>
                    <span class="workout-stat-value">${translateFeeling(workout.feeling)}</span>
                </div>
            </div>
            ${workout.notes ? `<p style="margin-top: 10px;"><em>${workout.notes}</em></p>` : ''}
            ${showAnalyzeBtn ? `<button class="btn btn-analyze" onclick="analyzeWorkout(${workout.id})">Analizar con IA</button>` : ''}
            <div id="analysis-${workout.id}"></div>
        </div>
    `;
}

// Manejar envío de nuevo workout
async function handleWorkoutSubmit(e) {
    e.preventDefault();
    
    const workoutData = {
        user_id: currentUser ? currentUser.id : 1,
        date: document.getElementById('date').value,
        type: document.getElementById('type').value,
        distance: parseFloat(document.getElementById('distance').value),
        duration: parseInt(document.getElementById('duration').value),
        avg_pace: document.getElementById('avg-pace').value,
        avg_heart_rate: parseInt(document.getElementById('avg-heart-rate').value) || 0,
        avg_power: parseInt(document.getElementById('avg-power').value) || 0,
        cadence: parseInt(document.getElementById('cadence').value) || 0,
        elevation_gain: parseInt(document.getElementById('elevation-gain').value) || 0,
        calories: parseInt(document.getElementById('calories').value) || 0,
        feeling: document.getElementById('feeling').value,
        notes: document.getElementById('notes').value
    };
    
    try {
        showLoading(true);
        
        // Primero analizar con OpenAI
        const response = await fetch(`${API_URL}/workout-analysis-form`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(workoutData)
        });
        
        if (response.ok) {
            const result = await response.json();
            const resultDiv = document.getElementById('workout-result');
            
            // Mostrar análisis con markdown y opción de guardar
            resultDiv.innerHTML = `
                <div class="analysis-section">
                    <h3>📊 Análisis del Entreno</h3>
                    <div class="markdown-content">${marked.parse(result.analysis)}</div>
                    
                    <div style="margin-top: 20px; padding: 15px; background: var(--bg-color); border-radius: 8px;">
                        <h4>✅ ¿Guardar este entreno en el historial?</h4>
                        <button onclick="saveWorkoutFromForm(${JSON.stringify(workoutData).replace(/"/g, '&quot;')})" class="btn btn-primary" style="margin-top: 15px;">
                            💾 Guardar en el Historial
                        </button>
                    </div>
                    
                    <!-- Chat de conversación -->
                    <div class="chat-container" style="margin-top: 20px;">
                        <div class="chat-messages" id="form-analysis-messages"></div>
                        <div class="chat-input-group">
                            <input type="text" id="form-analysis-question" placeholder="Pregunta algo sobre este análisis..." class="chat-input">
                            <button onclick="askAboutFormAnalysis()" class="btn btn-secondary" style="background: var(--primary-color); color: white;">Enviar</button>
                        </div>
                    </div>
                </div>
            `;
            
            // Añadir el análisis como primer mensaje
            addChatMessage('assistant', result.analysis, 'form-analysis-messages');
            
            // Permitir enviar con Enter
            const input = document.getElementById('form-analysis-question');
            input.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    askAboutFormAnalysis();
                }
            });
            
            resultDiv.style.display = 'block';
            resultDiv.scrollIntoView({ behavior: 'smooth' });
        } else {
            showResult('workout-result', '❌ Error al analizar el entreno', 'error');
        }
    } catch (error) {
        console.error('Error:', error);
        showResult('workout-result', '❌ Error de conexión', 'error');
    } finally {
        showLoading(false);
    }
}

// Guardar workout desde formulario después del análisis
async function saveWorkoutFromForm(workoutData) {
    try {
        showLoading(true);
        
        // Convertir fecha al formato RFC3339 que espera Go
        if (workoutData.date && !workoutData.date.includes('Z') && !workoutData.date.includes('+')) {
            // Si la fecha es del tipo "2025-12-02T16:19", añadir segundos y zona horaria
            workoutData.date = workoutData.date + ':00Z';
        }
        
        const response = await fetch(`${API_URL}/workouts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(workoutData)
        });
        
        if (response.ok) {
            alert('✅ ¡Entreno guardado en el historial exitosamente!');
            document.getElementById('workout-form').reset();
            
            // Configurar fecha actual nuevamente
            const now = new Date();
            now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
            document.getElementById('date').value = now.toISOString().slice(0, 16);
            
            await loadWorkouts();
            
            // Cambiar a la pestaña de workouts
            showTab('workouts');
        } else {
            alert('❌ Error al guardar el entreno');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('❌ Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Preguntar sobre el análisis del formulario
async function askAboutFormAnalysis() {
    const input = document.getElementById('form-analysis-question');
    const question = input.value.trim();
    
    if (!question) return;
    
    try {
        // Añadir pregunta del usuario al chat
        addChatMessage('user', question, 'form-analysis-messages');
        input.value = '';
        
        showLoading(true);
        
        // Enviar pregunta al backend
        const response = await fetch(`${API_URL}/workout-analysis-form`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ question: question })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Añadir respuesta del asistente al chat
            addChatMessage('assistant', result.analysis, 'form-analysis-messages');
            
            // Scroll al último mensaje
            const messagesDiv = document.getElementById('form-analysis-messages');
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        } else {
            alert('Error al procesar la pregunta');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Analizar workout con IA
async function analyzeWorkout(workoutId) {
    const analysisDiv = document.getElementById(`analysis-${workoutId}`);
    
    try {
        showLoading(true);
        const response = await fetch(`${API_URL}/workout-analysis`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ workout_id: workoutId })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Crear contenedor con análisis renderizado en markdown y chat
            analysisDiv.innerHTML = `
                <div class="analysis-section">
                    <h4>📊 Análisis con IA</h4>
                    <div class="markdown-content">${marked.parse(result.analysis)}</div>
                    
                    <!-- Chat de conversación sobre el análisis -->
                    <div class="chat-container" style="margin-top: 20px;">
                        <div class="chat-messages" id="analysis-messages-${workoutId}"></div>
                        <div class="chat-input-group">
                            <input type="text" id="analysis-question-${workoutId}" placeholder="Pregunta algo sobre este análisis..." class="chat-input">
                            <button onclick="askAboutAnalysis(${workoutId})" class="btn btn-secondary">Enviar</button>
                        </div>
                    </div>
                </div>
            `;
            
            // Añadir el análisis como primer mensaje
            addChatMessage('assistant', result.analysis, `analysis-messages-${workoutId}`);
            
            // Permitir enviar con Enter
            const input = document.getElementById(`analysis-question-${workoutId}`);
            input.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    askAboutAnalysis(workoutId);
                }
            });
            
            // Scroll al análisis
            analysisDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        } else {
            analysisDiv.innerHTML = '<p style="color: red;">Error al analizar el entreno</p>';
        }
    } catch (error) {
        console.error('Error:', error);
        analysisDiv.innerHTML = '<p style="color: red;">Error de conexión</p>';
    } finally {
        showLoading(false);
    }
}

// Preguntar sobre el análisis
async function askAboutAnalysis(workoutId) {
    const input = document.getElementById(`analysis-question-${workoutId}`);
    const question = input.value.trim();
    
    if (!question) return;
    
    try {
        // Añadir pregunta del usuario al chat
        addChatMessage('user', question, `analysis-messages-${workoutId}`);
        input.value = '';
        
        showLoading(true);
        
        // Enviar pregunta al backend (usando el mismo endpoint con la pregunta)
        const response = await fetch(`${API_URL}/workout-analysis`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                workout_id: workoutId,
                question: question 
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Añadir respuesta del asistente al chat
            addChatMessage('assistant', result.analysis, `analysis-messages-${workoutId}`);
            
            // Scroll al último mensaje
            const messagesDiv = document.getElementById(`analysis-messages-${workoutId}`);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        } else {
            alert('Error al procesar la pregunta');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Manejar envío de plan de entrenamiento
async function handlePlanSubmit(e) {
    e.preventDefault();
    
    const goal = document.getElementById('goal').value;
    
    try {
        showLoading(true);
        const response = await fetch(`${API_URL}/training-plan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                user_id: currentUser ? currentUser.id : 1,
                goal: goal 
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            document.getElementById('plan-content').textContent = result.plan;
            document.getElementById('plan-result').style.display = 'block';
        } else {
            alert('Error al generar el plan de entrenamiento');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Manejar envío de informe de progreso
async function handleReportSubmit(e) {
    e.preventDefault();
    
    const periodStart = document.getElementById('period-start').value;
    const periodEnd = document.getElementById('period-end').value;
    
    try {
        showLoading(true);
        const response = await fetch(`${API_URL}/progress-report`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                user_id: currentUser ? currentUser.id : 1,
                period_start: periodStart,
                period_end: periodEnd
            })
        });
        
        if (response.ok) {
            const result = await response.json();
            document.getElementById('report-content').textContent = result.report;
            document.getElementById('report-result').style.display = 'block';
        } else {
            alert('Error al generar el informe');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Error de conexión');
    } finally {
        showLoading(false);
    }
}

// Utilidades
function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'flex' : 'none';
}

function showResult(elementId, message, type) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.style.display = 'block';
    element.style.color = type === 'success' ? 'green' : 'red';
    
    setTimeout(() => {
        element.style.display = 'none';
    }, 5000);
}

function translateWorkoutType(type) {
    const types = {
        'easy': 'Carrera Suave',
        'interval': 'Intervalos',
        'tempo': 'Tempo',
        'long_run': 'Carrera Larga'
    };
    return types[type] || type;
}

function translateFeeling(feeling) {
    const feelings = {
        'great': '🌟 Excelente',
        'good': '😊 Bien',
        'ok': '😐 Normal',
        'tired': '😓 Cansado',
        'exhausted': '😫 Agotado'
    };
    return feelings[feeling] || feeling;
}

function translateFitnessLevel(level) {
    const levels = {
        'beginner': 'Principiante',
        'intermediate': 'Intermedio',
        'advanced': 'Avanzado'
    };
    return levels[level] || level;
}

// ==================== NUEVAS FUNCIONES PARA MÉTRICAS COMPARATIVAS ====================

// Estado del período seleccionado
let currentPeriod = 'week';

// Actualizar dashboard con período específico
function updateDashboardPeriod(period) {
    currentPeriod = period;
    
    // Actualizar botones
    document.querySelectorAll('.period-btn').forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
    
    // Calcular y mostrar métricas
    updateComparativeMetrics();
    updateWeeklyChart();
}

// Calcular métricas comparativas
function updateComparativeMetrics() {
    const now = new Date();
    let currentPeriodWorkouts = [];
    let previousPeriodWorkouts = [];
    
    if (currentPeriod === 'week') {
        // Esta semana (lunes a domingo)
        const startOfWeek = getStartOfWeek(now);
        const endOfWeek = new Date(startOfWeek);
        endOfWeek.setDate(endOfWeek.getDate() + 7);
        
        currentPeriodWorkouts = allWorkouts.filter(w => {
            const date = new Date(w.date);
            return date >= startOfWeek && date < endOfWeek;
        });
        
        // Semana anterior
        const prevStartOfWeek = new Date(startOfWeek);
        prevStartOfWeek.setDate(prevStartOfWeek.getDate() - 7);
        
        previousPeriodWorkouts = allWorkouts.filter(w => {
            const date = new Date(w.date);
            return date >= prevStartOfWeek && date < startOfWeek;
        });
        
    } else if (currentPeriod === 'month') {
        // Este mes
        const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
        const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
        
        currentPeriodWorkouts = allWorkouts.filter(w => {
            const date = new Date(w.date);
            return date >= startOfMonth && date <= endOfMonth;
        });
        
        // Mes anterior
        const prevStartOfMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1);
        const prevEndOfMonth = new Date(now.getFullYear(), now.getMonth(), 0);
        
        previousPeriodWorkouts = allWorkouts.filter(w => {
            const date = new Date(w.date);
            return date >= prevStartOfMonth && date <= prevEndOfMonth;
        });
        
    } else {
        // Total (sin comparativa)
        currentPeriodWorkouts = allWorkouts;
        previousPeriodWorkouts = [];
    }
    
    // Calcular métricas
    const currentStats = calculateStats(currentPeriodWorkouts);
    const previousStats = calculateStats(previousPeriodWorkouts);
    
    // Mostrar métricas
    document.getElementById('period-workouts').textContent = currentStats.count;
    document.getElementById('period-distance').textContent = `${currentStats.distance.toFixed(1)} km`;
    document.getElementById('period-pace').textContent = currentStats.avgPace || '--:--';
    document.getElementById('period-hr').textContent = currentStats.avgHR ? `${currentStats.avgHR} bpm` : '-- bpm';
    
    // Mostrar tendencias (solo si hay período anterior)
    if (previousPeriodWorkouts.length > 0 && currentPeriod !== 'all') {
        updateTrend('workouts-trend', currentStats.count, previousStats.count);
        updateTrend('distance-trend', currentStats.distance, previousStats.distance, 'km');
        updateTrendPace('pace-trend', currentStats.avgPaceSeconds, previousStats.avgPaceSeconds);
        updateTrend('hr-trend', currentStats.avgHR, previousStats.avgHR, 'bpm');
    } else {
        // Sin tendencias
        ['workouts-trend', 'distance-trend', 'pace-trend', 'hr-trend'].forEach(id => {
            document.getElementById(id).textContent = '';
            document.getElementById(id).className = 'trend';
        });
    }
}

// Calcular estadísticas de un conjunto de workouts
function calculateStats(workouts) {
    if (workouts.length === 0) {
        return { count: 0, distance: 0, avgPace: null, avgPaceSeconds: 0, avgHR: 0 };
    }
    
    const count = workouts.length;
    const distance = workouts.reduce((sum, w) => sum + (w.distance || 0), 0);
    
    // Calcular ritmo promedio (convertir a segundos)
    const validPaces = workouts.filter(w => w.avg_pace && w.avg_pace !== '--:--');
    let avgPaceSeconds = 0;
    if (validPaces.length > 0) {
        const totalSeconds = validPaces.reduce((sum, w) => {
            const [min, sec] = w.avg_pace.split(':').map(Number);
            return sum + (min * 60 + sec);
        }, 0);
        avgPaceSeconds = totalSeconds / validPaces.length;
    }
    
    const avgPace = avgPaceSeconds > 0 ? formatPace(avgPaceSeconds) : null;
    
    // Calcular FC promedio
    const validHR = workouts.filter(w => w.avg_heart_rate && w.avg_heart_rate > 0);
    const avgHR = validHR.length > 0 
        ? Math.round(validHR.reduce((sum, w) => sum + w.avg_heart_rate, 0) / validHR.length)
        : 0;
    
    return { count, distance, avgPace, avgPaceSeconds, avgHR };
}

// Formatear ritmo de segundos a MM:SS
function formatPace(seconds) {
    const min = Math.floor(seconds / 60);
    const sec = Math.round(seconds % 60);
    return `${min}:${sec.toString().padStart(2, '0')}`;
}

// Actualizar indicador de tendencia
function updateTrend(elementId, current, previous, unit = '') {
    const element = document.getElementById(elementId);
    
    if (previous === 0) {
        element.textContent = '';
        element.className = 'trend';
        return;
    }
    
    const diff = current - previous;
    const percentChange = ((diff / previous) * 100).toFixed(1);
    
    element.className = 'trend';
    
    if (Math.abs(percentChange) < 2) {
        element.classList.add('neutral');
        element.textContent = 'Similar';
    } else if (diff > 0) {
        element.classList.add('up');
        element.textContent = `+${Math.abs(percentChange)}%`;
    } else {
        element.classList.add('down');
        element.textContent = `${Math.abs(percentChange)}%`;
    }
}

// Actualizar tendencia de ritmo (inverso: menor es mejor)
function updateTrendPace(elementId, currentSeconds, previousSeconds) {
    const element = document.getElementById(elementId);
    
    if (previousSeconds === 0 || currentSeconds === 0) {
        element.textContent = '';
        element.className = 'trend';
        return;
    }
    
    const diff = currentSeconds - previousSeconds;
    const percentChange = ((Math.abs(diff) / previousSeconds) * 100).toFixed(1);
    
    element.className = 'trend';
    
    if (Math.abs(diff) < 2) {
        element.classList.add('neutral');
        element.textContent = 'Similar';
    } else if (diff < 0) {
        // Ritmo más rápido (menos segundos) = mejor
        element.classList.add('up');
        element.textContent = `${percentChange}% más rápido`;
    } else {
        // Ritmo más lento (más segundos) = peor
        element.classList.add('down');
        element.textContent = `${percentChange}% más lento`;
    }
}

// Obtener inicio de semana (lunes)
function getStartOfWeek(date) {
    const d = new Date(date);
    const day = d.getDay();
    const diff = d.getDate() - day + (day === 0 ? -6 : 1);
    return new Date(d.setDate(diff));
}

// Actualizar gráfica semanal
function updateWeeklyChart() {
    const now = new Date();
    const startOfWeek = getStartOfWeek(now);
    
    // Crear array de 7 días con distancias
    const weekData = [];
    const dayNames = ['L', 'M', 'X', 'J', 'V', 'S', 'D'];
    
    for (let i = 0; i < 7; i++) {
        const day = new Date(startOfWeek);
        day.setDate(day.getDate() + i);
        
        const dayWorkouts = allWorkouts.filter(w => {
            const workoutDate = new Date(w.date);
            return workoutDate.toDateString() === day.toDateString();
        });
        
        const distance = dayWorkouts.reduce((sum, w) => sum + (w.distance || 0), 0);
        
        weekData.push({
            label: dayNames[i],
            value: distance,
            date: day
        });
    }
    
    // Renderizar gráfica
    const chartContainer = document.getElementById('weekly-chart');
    const maxDistance = Math.max(...weekData.map(d => d.value), 1);
    
    chartContainer.innerHTML = weekData.map(day => {
        const heightPercent = (day.value / maxDistance) * 100;
        return `
            <div class="chart-bar" style="height: ${heightPercent}%" title="${day.label}: ${day.value.toFixed(1)} km">
                <span class="chart-bar-value">${day.value > 0 ? day.value.toFixed(1) : ''}</span>
                <span class="chart-bar-label">${day.label}</span>
            </div>
        `;
    }).join('');
}

// Sobrescribir updateDashboardStats para usar el nuevo sistema
const originalUpdateDashboardStats = updateDashboardStats;
updateDashboardStats = function() {
    updateComparativeMetrics();
    updateWeeklyChart();
};

// ==================== INTEGRACIÓN CON STRAVA ====================

// Verificar estado de conexión con Strava
async function checkStravaStatus() {
    try {
        const response = await fetch(`${API_URL}/strava/status`);
        const data = await response.json();
        
        if (data.connected) {
            document.getElementById('strava-connected').style.display = 'block';
            document.getElementById('strava-disconnected').style.display = 'none';
            
            if (data.last_sync) {
                const lastSync = new Date(data.last_sync);
                document.getElementById('strava-last-sync').textContent = lastSync.toLocaleString('es-ES');
            }
        } else {
            document.getElementById('strava-connected').style.display = 'none';
            document.getElementById('strava-disconnected').style.display = 'block';
        }
    } catch (error) {
        console.error('Error verificando estado de Strava:', error);
    }
}

// Conectar con Strava
function connectStrava() {
    window.location.href = `${API_URL}/strava/auth`;
}

// Sincronizar entrenos de Strava
async function syncStrava() {
    try {
        showLoading(true);
        
        const response = await fetch(`${API_URL}/strava/sync`, {
            method: 'POST'
        });
        
        if (response.ok) {
            const result = await response.json();
            
            if (result.imported > 0) {
                alert(`✅ Se importaron ${result.imported} nuevos entrenamientos de Strava!`);
                await loadWorkouts(); // Recargar lista de workouts
                checkStravaStatus(); // Actualizar estado
            } else {
                alert('ℹ️ No hay entrenamientos nuevos para sincronizar');
            }
        } else {
            const error = await response.text();
            alert('❌ Error al sincronizar: ' + error);
        }
    } catch (error) {
        console.error('Error sincronizando Strava:', error);
        alert('❌ Error de conexión al sincronizar con Strava');
    } finally {
        showLoading(false);
    }
}

