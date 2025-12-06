User Manual
===========

6 Dec 2025 by White Water

Overview
--------

This document is a concise user manual for the ERP project. It focuses on the current development status and simple instructions to run and use the system.

Development Progress (summary)
------------------------------

- Core backend implemented (Go, Gin, GORM) with model, repository, service and HTTP handler layers.
- Authentication and role-based access (donor, volunteer, employee) implemented.
- Dynamic CURD functions over UI implemented — metadata-driven forms and tables.
- Role-specific front-end pages created: donor, volunteer, employee dashboards.
- Reporting/charting support added (Chart.js) with backend aggregation endpoints.
- TODOs: all business logic; clean up for DBMS UI.

Quick Start (3 ways)
--------------------

1. Visit the webpage held by our team

- [https://ecf-erp.bugnomad.top](https://ecf-erp.bugnomad.top)
- This might not be stable since deployed on personal on-permis system.

2. Run the pre build binaries

- Run the executale in your shell
- Windows `/backend/bin/erp-backend.exe`
- Linux `/backend/bin/erp-backend-linux-amd64`
- No path should be required since relative path applied
- configration file should be an `.env` file. Refer to the file in `/backend/internal/config` for detail.

3. Build and run

- Not reocomended
- Source code available: zipfile (with test data); github https://github.com/chkviii/MIS-for-ECF.git (code only)

Basic Usage (for each role)
---------------------------

- Landing at welcome page, login or register

  - Note: register as employee won't grant access to the ERP system. Change `status` field in `users` table from `pending` to `active` either by directly minipulating the database, or login as `EmpTestUser1` and use `System Admin`.
- Employee (admin):

  - Existing account (username and password): `EmpTestUser1`
  - Access the Dadabase management interface through `DBMS` after login.
  - Select a visual finantial report from the sidebar to view.
  - Manage syetem user through `System Admin`.
  - Projects management and Home page not implemented yet.
- Donor:

  - Existing account (username and password): `DonTestUser1`
  - Donor Dashboard and Donation History is available.
- Volunteer:

  - Existing account (username and password): `VolTestUser1`
  - Not functioning
- Logging in with given accounts is recomended. Injected data can showcase some functions, otherwise manual input of original data is required, especially for charts.

Reporting & Charts
------------------

- Charts are available on the employee dashboard and the report page.
- Most of the controls should be working but

Where to find things
--------------------

- Source code: zip file (with test data); Github Repository [https://github.com/chkviii/MIS-for-ECF.git](https://github.com/chkviii/MIS-for-ECF.git) (code only)
- Backend code: `backend/internal/` (models, repo, services, handlers)
- Frontend templates: `frontend/templates/`
- Frontend JS/CSS: `frontend/static/js/` and `frontend/static/css/`

Contact
-------

If you need help with running or using the system, please contact me for technical support.
