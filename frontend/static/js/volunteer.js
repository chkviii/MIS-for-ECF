import { apiRequestSimple, checkAuthSimple, logoutSimple, API_BASE_URL } from './site-api.js';

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuthSimple()) return;

    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const userId = user.user_id;

    // Resolve volunteer record
    let volunteerId = null;
    try {
        const resp = await apiRequestSimple(`${API_BASE_URL}/volunteers/search?query=` + encodeURIComponent(JSON.stringify({ user_id: userId })));
        const data = await resp.json();
        if (data && data.data && data.data.length > 0) {
            volunteerId = data.data[0].id;
        }
    } catch (e) {
        console.error('Failed to fetch volunteer record', e);
    }

    // Load projects
    try {
        const resp = await apiRequestSimple(`${API_BASE_URL}/projects`);
        const result = await resp.json();
        const list = (result && result.data) ? result.data : [];
        const container = document.getElementById('projects');
        container.innerHTML = '';
        list.forEach(p => {
            const div = document.createElement('div');
            div.className = 'project';
            div.innerHTML = `
                <strong>${p.name}</strong>
                <div style="margin-top:6px;color:#666">${p.description || ''}</div>
                <div style="margin-top:8px"><button data-id="${p.id}" class="apply-btn">Apply</button></div>
            `;
            container.appendChild(div);
        });

        container.addEventListener('click', async (ev) => {
            const btn = ev.target.closest('.apply-btn');
            if (!btn) return;
            const projectId = Number(btn.getAttribute('data-id'));
            if (!volunteerId) {
                alert('Volunteer profile not found. Please complete your volunteer profile first.');
                return;
            }
            try {
                const body = { volunteer_id: volunteerId, project_id: projectId, role: 'applicant', status: 'pending' };
                const r = await apiRequestSimple(`${API_BASE_URL}/volunteer-projects`, { method: 'POST', body: JSON.stringify(body) });
                if (r.ok) {
                    alert('Applied successfully');
                } else {
                    const err = await r.json();
                    alert('Failed to apply: ' + (err.error || JSON.stringify(err)));
                }
            } catch (e) {
                console.error('Apply failed', e);
                alert('Apply failed');
            }
        });

    } catch (e) {
        console.error('Failed to load projects', e);
    }

    // Load volunteer history (volunteer-projects)
    try {
        let historyHtml = '';
        let totalHours = 0;
        if (volunteerId) {
            const resp = await apiRequestSimple(`${API_BASE_URL}/volunteer-projects/search?query=` + encodeURIComponent(JSON.stringify({ volunteer_id: volunteerId })));
            const res = await resp.json();
            const items = res && res.data ? res.data : [];
            items.forEach(it => {
                historyHtml += `<div style="padding:8px;border-bottom:1px solid #eee">Project: ${it.project ? it.project.name : it.project_id} — Role: ${it.role || ''} — Status: ${it.status || ''}</div>`;
            });

            // Sum hours from schedules table
            const sch = await apiRequestSimple(`${API_BASE_URL}/schedules/search?query=` + encodeURIComponent(JSON.stringify({ person_type: 'volunteer', person_id: volunteerId })));
            const schRes = await sch.json();
            const schList = schRes && schRes.data ? schRes.data : [];
            schList.forEach(s => { totalHours += Number(s.hours_worked || 0); });
        } else {
            historyHtml = '<div>No volunteer profile found.</div>';
        }

        document.getElementById('history').innerHTML = historyHtml;
        document.getElementById('total-hours').textContent = 'Total hours: ' + totalHours;
    } catch (e) {
        console.error('Failed to load volunteer history', e);
    }
});

window.logoutSimple = logoutSimple;
