import {apiRequest, showToast} from '/static/js/site-api.js';

let ProjectPromise = null;
// const projects = {};

document.addEventListener('DOMContentLoaded', function () {
    // Initialize donor dashboard components
    initDonorDashboard();

    // Initialize donation chart
    initDonationChart();

    // Initialize donation details
    initDetails();
});

function initDonorDashboard() {
    const navDashboard = document.getElementById('dashboardnav');
    const navDonations = document.getElementById('donationsnav');
    const makeDon = document.getElementById('makedon');
    const secDashboard = document.getElementById('chart-section');
    const secDonations = document.getElementById('details-section');
    const secMakeDon = document.getElementById('new-donation-section');
    console.log('Donor dashboard initialized', navDashboard, navDonations, makeDon, secDashboard, secDonations, secMakeDon);

    secDashboard.style.display = 'block';
    secDonations.style.display = 'none';
    secMakeDon.style.display = 'none';

    navDashboard.addEventListener('click', function () {
        secDashboard.style.display = 'block';
        secDonations.style.display = 'none';
        secMakeDon.style.display = 'none';
    });

    navDonations.addEventListener('click', async function () {
        secDashboard.style.display = 'none';
        secDonations.style.display = 'block';
        secMakeDon.style.display = 'none';
    });

    makeDon.addEventListener('click', function () {
        secDashboard.style.display = 'none';
        secDonations.style.display = 'none';
        secMakeDon.style.display = 'block';
    });
}

let ReportModulePromise = null;

// Add report event listeners
async function initDonationChart() {
        try {
            // ensure module loaded (use cached promise if present)
            const mod = ReportModulePromise ? await ReportModulePromise : await import('/static/js/report.js');
            const ReportChart = mod && (mod.default || mod.ReportChart);
            if (!ReportChart) throw new Error('ReportChart class not found in module');

            // reuse a single instance to avoid duplicates; destroy previous if exists
            if (window._currentReportChart) {
                try { window._currentReportChart.destroy(); } catch(e){/*ignore*/ }
                window._currentReportChart = null;
            }

            const options = {
                urlbase: '/donor/charts',
                chartID: 'donations',
            };

            // console.log('Initializing report chart with options:', options);

            // instantiate and update
            window._currentReportChart = new ReportChart(options);
            await window._currentReportChart.update();
        } catch (err) {
            console.error('Failed to load or initialize report chart:', err);
            if (typeof showToast === 'function') showToast('Failed to load report: ' + (err.message || err), 'error');
        }

}

async function preloadProject() {
    if (!ProjectPromise) {
        const projectSelector = document.getElementById('dproject');
        ProjectPromise = apiRequest('/api/v1/donor/projects', { method: 'GET'})
            .then(async res => {
                const data = await res.json();
                if (!res.ok) throw new Error(`Failed to fetch projects: `, data.message);
                if (!data.projects || data.projects.length === 0) {
                    console.warn('No projects found for donor.');
                    return;
                }
                data.projects.forEach(proj => {
                    const option = document.createElement('option');
                    option.value = proj.id;
                    option.textContent = proj.name;
                    projectSelector.appendChild(option);
                }
                );
            })
            .catch(err => {
                console.error('Failed to preload projects:', err);
                // clear cache so future attempts can retry
                ProjectPromise = null;
            });
    }
}



async function initDetails() {
    const detailtable = document.getElementById('donation-details-table');
    if (!detailtable) return;

    const startInput = document.getElementById('dstartDate');
    const endInput = document.getElementById('dendDate');
    const projectInput = document.getElementById('dproject');
    const btnupdate = document.getElementById('dapplyBtn');
    const navbtn = document.getElementById('donationsnav');

    btnupdate?.addEventListener('click', async function () {
        await loadDetails(startInput.value, endInput.value, projectInput.value);
    });

    navbtn?.addEventListener('click', async function () {
        await loadDetails(startInput.value, endInput.value, projectInput.value);
    });
}

async function loadDetails(start, end, project) {

    const detailtable = document.getElementById('donation-details-table');

    let url = `/api/v1/donor/donations?`;
    if (project) url += `&project=${encodeURIComponent(project)}`;
    if (start) url += `&start=${encodeURIComponent(start)}`;
    if (end) url += `&end=${encodeURIComponent(end)}`;

    try {

        await preloadProject();

        const res = await apiRequest(url, { method: 'GET', cache: 'no-store' });
        const data = await res.json();
        if (!res.ok) throw new Error(`Failed to fetch donation details: `, data.message);

        // Populate table
        const tbody = detailtable.querySelector('tbody');
        tbody.innerHTML = '';
        if (data.details && data.details.length > 0) {
            data.details.forEach(donation => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${donation.id}</td>
                    <td>${donation.date}</td>
                    <td>${donation.project}</td>
                    <td>${donation.amount}</td>
                    <td>${donation.method}</td>
                `;
                tbody.appendChild(tr);
            });
        }
    } catch (err) {
        console.error('Error loading donation details:', err);
        showToast('Error loading donation details: ' + (err.message || err), 'error');
    }
}





        



