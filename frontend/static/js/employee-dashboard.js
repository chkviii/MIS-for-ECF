// ERP Management System JavaScript
const API_BASE_URL = '/api/v1';

let ReportModulePromise = null;
let chart = null;

// Add logout button event on page load
document.addEventListener('DOMContentLoaded', async function() {

    isAnyReportNavItemActive();
    addReportNavEventListeners();
});

// Find if any report nav item in sidebar is active
function isAnyReportNavItemActive() {
    const reportContainer = document.getElementById('report-container');
    if (reportContainer) reportContainer.style.display = 'none';

    const reportNavItems = document.querySelectorAll('.nav-item[data-report]');
    if (!reportNavItems || reportNavItems.length === 0) {
        // no report nav items -> keep hidden
        return;
    }

    // Preload report.js asynchronously and cache the promise
    if (!ReportModulePromise) {
        ReportModulePromise = import('/static/js/report.js')
            .catch(err => {
                console.error('Failed to preload report module:', err);
                // clear cache so future attempts can retry
                ReportModulePromise = null;
            });
    }
}

// Add report event listeners
function addReportNavEventListeners() {
    const reportNavItems = document.querySelectorAll('.nav-item[data-report]');
    // console.log('Report nav items:', reportNavItems);
    if (!reportNavItems || reportNavItems.length === 0) return;

    reportNavItems.forEach(item => {
        item.addEventListener('click', async function(e) {
            const mainContent = document.getElementById('main-content');
            for (const child of mainContent.children) {
                child.style.display = 'none';
            }

            const headertitle = document.getElementById('page-title');
            if (headertitle) {
                headertitle.textContent = item.textContent || 'Report';
            }

            const reportContainer = document.getElementById('report-container');
            if (reportContainer) reportContainer.style.display = 'block';

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

                options = {
                    urlbase: '/fin/charts',
                    chartID: item.getAttribute('data-report'),
                };

                // console.log('Initializing report chart with options:', options);

                // instantiate and update
                window._currentReportChart = new ReportChart(options);
                await window._currentReportChart.update();
            } catch (err) {
                console.error('Failed to load or initialize report chart:', err);
                if (typeof showToast === 'function') showToast('无法加载报表: ' + (err.message || err), 'error');
            }
        });
    });
}