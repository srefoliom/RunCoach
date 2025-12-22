// Workout Detail Page Logic
const API_URL = 'http://localhost:8080/api';

// Get workout ID from URL
function getWorkoutId() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

// Format duration (seconds to HH:MM:SS)
function formatDuration(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = Math.floor(seconds % 60);
    
    if (hours > 0) {
        return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    }
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
}

// Format pace (min/km)
function formatPace(metersPerSecond) {
    if (!metersPerSecond) return '-';
    const minutesPerKm = 1000 / (metersPerSecond * 60);
    const minutes = Math.floor(minutesPerKm);
    const seconds = Math.floor((minutesPerKm - minutes) * 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

// Format date
function formatDate(dateString) {
    const options = { 
        weekday: 'long', 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };
    return new Date(dateString).toLocaleDateString('es-ES', options);
}

// Render main stats
function renderMainStats(workout) {
    const statsGrid = document.getElementById('main-stats');
    const stats = [
        { label: '📏 Distancia', value: (workout.distance / 1000).toFixed(2), unit: 'km' },
        { label: '⏱️ Duración', value: formatDuration(workout.moving_time || workout.elapsed_time), unit: '' },
        { label: '⚡ Ritmo', value: formatPace(workout.average_speed), unit: 'min/km' },
        { label: '❤️ FC Media', value: workout.average_heartrate ? Math.round(workout.average_heartrate) : '-', unit: 'bpm' },
        { label: '💪 FC Máx', value: workout.max_heartrate ? Math.round(workout.max_heartrate) : '-', unit: 'bpm' },
        { label: '⚡ Potencia Media', value: workout.average_watts ? Math.round(workout.average_watts) : '-', unit: 'W' },
        { label: '🔥 Potencia Máx', value: workout.max_watts ? Math.round(workout.max_watts) : '-', unit: 'W' },
        { label: '👟 Cadencia', value: workout.average_cadence ? Math.round(workout.average_cadence * 2) : '-', unit: 'spm' },
        { label: '⛰️ Elevación', value: workout.total_elevation_gain ? Math.round(workout.total_elevation_gain) : '-', unit: 'm' },
        { label: '🔥 Calorías', value: workout.calories || '-', unit: 'kcal' },
        { label: '😊 Sensaciones', value: workout.perceived_exertion || '-', unit: '/10' },
        { label: '💯 Suffer Score', value: workout.suffer_score || '-', unit: '' }
    ];
    
    statsGrid.innerHTML = stats.map(stat => `
        <div class="stat-item">
            <span class="stat-label">${stat.label}</span>
            <span class="stat-value">${stat.value} <span class="stat-unit">${stat.unit}</span></span>
        </div>
    `).join('');
}

// Render best efforts
function renderBestEfforts(efforts) {
    if (!efforts || efforts.length === 0) {
        document.getElementById('best-efforts-card').style.display = 'none';
        return;
    }
    
    const effortsList = document.getElementById('best-efforts');
    effortsList.innerHTML = efforts.map(effort => {
        const prBadge = effort.pr_rank ? `<span class="effort-pr">PR #${effort.pr_rank}</span>` : '';
        return `
            <div class="effort-item">
                <span class="effort-name">${effort.name}</span>
                <span>
                    <span class="effort-time">${formatDuration(effort.elapsed_time)}</span>
                    ${prBadge}
                </span>
            </div>
        `;
    }).join('');
}

// Render map
function renderMap(mapData) {
    if (!mapData || !mapData.summary_polyline) {
        document.getElementById('map-card').style.display = 'none';
        return;
    }
    
    document.getElementById('map-card').style.display = 'block';
    
    // Decode polyline using @mapbox/polyline library
    const coordinates = polyline.decode(mapData.summary_polyline);
    
    // Wait for container to be visible before creating map
    setTimeout(() => {
        // Create map
        const map = L.map('map').setView([coordinates[0][0], coordinates[0][1]], 13);
        
        // Add tile layer (OpenStreetMap)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '© OpenStreetMap contributors'
        }).addTo(map);
        
        // Add route polyline
        L.polyline(coordinates, {
            color: '#00d4aa',
            weight: 4,
            opacity: 0.8
        }).addTo(map);
        
        // Add start marker
        L.marker(coordinates[0], {
            icon: L.icon({
                iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
                shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
                iconSize: [25, 41],
                iconAnchor: [12, 41]
            })
        }).addTo(map).bindPopup('Inicio');
        
        // Add end marker
        const lastPoint = coordinates[coordinates.length - 1];
        L.marker(lastPoint, {
            icon: L.icon({
                iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
                shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
                iconSize: [25, 41],
                iconAnchor: [12, 41]
            })
        }).addTo(map).bindPopup('Fin');
        
        // Fit bounds
        map.fitBounds(L.latLngBounds(coordinates));
        
        // Force map to recalculate size after container is visible
        map.invalidateSize();
    }, 100);
}

// Render elevation chart
function renderElevationChart(splits) {
    if (!splits || splits.length === 0) {
        document.getElementById('elevation-card').style.display = 'none';
        return;
    }
    
    document.getElementById('elevation-card').style.display = 'block';
    
    const ctx = document.getElementById('elevation-chart').getContext('2d');
    const labels = splits.map((_, i) => `Km ${i + 1}`);
    const elevations = splits.map(split => split.elevation_difference || 0);
    
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Diferencia de Elevación (m)',
                data: elevations,
                borderColor: '#00d4aa',
                backgroundColor: 'rgba(0, 212, 170, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    labels: { color: '#e0e0e0' }
                }
            },
            scales: {
                y: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { color: '#e0e0e0' }
                },
                x: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { color: '#e0e0e0' }
                }
            }
        }
    });
}

// Render pace chart
function renderPaceChart(splits) {
    if (!splits || splits.length === 0) {
        document.getElementById('pace-card').style.display = 'none';
        return;
    }
    
    document.getElementById('pace-card').style.display = 'block';
    
    const ctx = document.getElementById('pace-chart').getContext('2d');
    const labels = splits.map((_, i) => `Km ${i + 1}`);
    const paces = splits.map(split => {
        const pace = 1000 / (split.average_speed * 60); // min/km
        return pace;
    });
    
    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Ritmo (min/km)',
                data: paces,
                backgroundColor: '#00d4aa',
                borderColor: '#00d4aa',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    labels: { color: '#e0e0e0' }
                }
            },
            scales: {
                y: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { 
                        color: '#e0e0e0',
                        callback: function(value) {
                            const minutes = Math.floor(value);
                            const seconds = Math.floor((value - minutes) * 60);
                            return `${minutes}:${seconds.toString().padStart(2, '0')}`;
                        }
                    }
                },
                x: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { color: '#e0e0e0' }
                }
            }
        }
    });
}

// Render HR chart
function renderHRChart(splits) {
    if (!splits || splits.length === 0 || !splits[0].average_heartrate) {
        document.getElementById('hr-card').style.display = 'none';
        return;
    }
    
    document.getElementById('hr-card').style.display = 'block';
    
    const ctx = document.getElementById('hr-chart').getContext('2d');
    const labels = splits.map((_, i) => `Km ${i + 1}`);
    const heartrates = splits.map(split => split.average_heartrate || 0);
    
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: 'Frecuencia Cardíaca (bpm)',
                data: heartrates,
                borderColor: '#ff4d4d',
                backgroundColor: 'rgba(255, 77, 77, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    labels: { color: '#e0e0e0' }
                }
            },
            scales: {
                y: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { color: '#e0e0e0' }
                },
                x: {
                    grid: { color: 'rgba(255, 255, 255, 0.1)' },
                    ticks: { color: '#e0e0e0' }
                }
            }
        }
    });
}

// Render splits table
function renderSplitsTable(splits) {
    if (!splits || splits.length === 0) {
        document.getElementById('splits-card').style.display = 'none';
        return;
    }
    
    document.getElementById('splits-card').style.display = 'block';
    
    const tbody = document.querySelector('.splits-table tbody');
    
    // Find fastest and slowest
    const speeds = splits.map(s => s.average_speed).filter(s => s > 0);
    const fastest = Math.max(...speeds);
    const slowest = Math.min(...speeds);
    
    tbody.innerHTML = splits.map((split, i) => {
        let rowClass = '';
        if (split.average_speed === fastest) rowClass = 'split-fastest';
        if (split.average_speed === slowest) rowClass = 'split-slowest';
        
        return `
            <tr class="${rowClass}">
                <td>${i + 1}</td>
                <td>${formatDuration(split.elapsed_time)}</td>
                <td>${formatPace(split.average_speed)}</td>
                <td>${split.average_heartrate ? Math.round(split.average_heartrate) : '-'}</td>
                <td>${split.elevation_difference ? split.elevation_difference.toFixed(1) : '-'}</td>
            </tr>
        `;
    }).join('');
}

// Render segments
function renderSegments(segments) {
    if (!segments || segments.length === 0) {
        document.getElementById('segments-card').style.display = 'none';
        return;
    }
    
    document.getElementById('segments-card').style.display = 'block';
    
    const segmentsList = document.getElementById('segments-list');
    segmentsList.innerHTML = segments.map(segment => {
        const komBadge = segment.kom_rank ? `<span class="segment-kom">KOM #${segment.kom_rank}</span>` : '';
        const prBadge = segment.pr_rank ? `<span class="segment-pr">PR #${segment.pr_rank}</span>` : '';
        
        return `
            <div class="segment-item">
                <span class="segment-name">${segment.name} ${komBadge} ${prBadge}</span>
                <div class="segment-stats">
                    <span class="segment-stat"><strong>Tiempo:</strong> ${formatDuration(segment.elapsed_time)}</span>
                    <span class="segment-stat"><strong>Distancia:</strong> ${(segment.distance / 1000).toFixed(2)} km</span>
                    ${segment.average_heartrate ? `<span class="segment-stat"><strong>FC Media:</strong> ${Math.round(segment.average_heartrate)} bpm</span>` : ''}
                </div>
            </div>
        `;
    }).join('');
}

// Render gear
function renderGear(gear) {
    if (!gear) {
        document.getElementById('gear-card').style.display = 'none';
        return;
    }
    
    const gearInfo = document.getElementById('gear-info');
    gearInfo.innerHTML = `
        <span class="gear-name">👟 ${gear.name}</span>
        <span class="gear-distance">Distancia acumulada: ${(gear.distance / 1000).toFixed(1)} km</span>
    `;
}

// Render achievements
function renderAchievements(workout) {
    const achievementsInfo = document.getElementById('achievements-info');
    const badges = [];
    
    if (workout.achievement_count > 0) {
        badges.push(`🏆 ${workout.achievement_count} logros`);
    }
    
    if (workout.pr_count > 0) {
        badges.push(`⭐ ${workout.pr_count} récords`);
    }
    
    if (workout.kudos_count > 0) {
        badges.push(`👏 ${workout.kudos_count} kudos`);
    }
    
    if (badges.length === 0) {
        document.getElementById('achievements-card').style.display = 'none';
        return;
    }
    
    achievementsInfo.innerHTML = badges.map(badge => 
        `<div class="achievement-badge">${badge}</div>`
    ).join('');
}

// Fetch workout detail
async function fetchWorkoutDetail(workoutId) {
    try {
        const token = localStorage.getItem('auth_token');
        const response = await fetch(`${API_URL}/workouts/${workoutId}/detail`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        
        if (!response.ok) {
            throw new Error('Error al cargar el detalle del entreno');
        }
        
        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        alert('Error al cargar el detalle del entreno');
        window.location.href = '/index.html';
    }
}

// Initialize page
async function init() {
    const workoutId = getWorkoutId();
    
    if (!workoutId) {
        alert('ID de entreno no válido');
        window.location.href = '/index.html';
        return;
    }
    
    try {
        // Fetch workout detail
        const workout = await fetchWorkoutDetail(workoutId);
        
        if (!workout) return;
        
        // Update header
        document.getElementById('workout-title').textContent = workout.name || 'Entreno';
        document.getElementById('workout-date').textContent = formatDate(workout.start_date);
        
        // Render sections
        renderMainStats(workout);
        renderBestEfforts(workout.best_efforts);
        renderGear(workout.gear);
        renderAchievements(workout);
        
        // Render map and charts
        if (workout.map && workout.map.summary_polyline) {
            renderMap(workout.map);
        } else {
            document.getElementById('map-card').style.display = 'none';
        }
        
        if (workout.splits_metric) {
            renderSplitsTable(workout.splits_metric);
            renderElevationChart(workout.splits_metric);
            renderPaceChart(workout.splits_metric);
            renderHRChart(workout.splits_metric);
        }
        
        if (workout.segment_efforts) {
            renderSegments(workout.segment_efforts);
        }
        
        // Hide loading, show content
        document.getElementById('loading').style.display = 'none';
        document.querySelector('.workout-detail-page').style.display = 'block';
    } catch (error) {
        console.error('Error initializing page:', error);
        document.getElementById('loading').style.display = 'none';
        alert('Error al cargar el detalle del entreno');
        window.location.href = '/index.html';
    }
}

// Initialize on load
document.addEventListener('DOMContentLoaded', init);
