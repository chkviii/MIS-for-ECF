import { apiRequestSimple, checkAuthSimple, logoutSimple, API_BASE_URL } from './site-api.js';

document.addEventListener('DOMContentLoaded', async () => {
    if (!checkAuthSimple()) return;

    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const userId = user.user_id;

    // Find donor record associated with this user
    let donorId = null;
    try {
        const resp = await apiRequestSimple(`${API_BASE_URL}/donors/search?query=` + encodeURIComponent(JSON.stringify({ user_id: userId })));
        const data = await resp.json();
        if (data && data.data && data.data.length > 0) {
            donorId = data.data[0].id;
        }
    } catch (e) {
        console.error('Failed to get donor record', e);
    }

    // Fetch all donations and filter by donor id
    try {
        const resp = await apiRequestSimple(`${API_BASE_URL}/donations`);
        const result = await resp.json();
        const list = (result && result.data) ? result.data : [];
        const filtered = donorId ? list.filter(d => d.donor_id === donorId) : [];

        const tbody = document.querySelector('#donation-table tbody');
        tbody.innerHTML = '';
        filtered.forEach((d, i) => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${i + 1}</td>
                <td>${new Date(d.donation_date).toLocaleDateString()}</td>
                <td>${d.amount}</td>
                <td>${d.donation_type || ''}</td>
                <td>${d.project ? (d.project.name || d.project.project_id) : ''}</td>
                <td>${d.notes || ''}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (e) {
        console.error('Failed to load donations', e);
    }
});

window.logoutSimple = logoutSimple;
