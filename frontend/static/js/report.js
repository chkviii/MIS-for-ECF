// report.js - ReportChart class (extracted from inline script)
import '/static/js/lib/chart.umd.min.js';
import '/static/js/lib/chartjs-adapter-date-fns.bundle.min.js';


class ReportChart{

    constructor(opts = {}){
        this.canvasId = opts.canvasId || 'mainChart';
        this.urlbase = opts.urlbase
        this.chartID = opts.chartID || 'mainChart';
        this.chartType = opts.chartType || 'chartType';
        this.startInputId = opts.startInputId || 'startDate';
        this.endInputId = opts.endInputId || 'endDate';
        this.applyBtnId = opts.applyBtnId || 'applyBtn';

        this.metricInputId = opts.metricInputId || 'metricInput';
        this.metricTotalId = opts.metricTotalId || 'metricTotal';
        this.averageId = opts.averageId || 'average';

        this._ctx = null;
        this._chart = null;
        this._listeners = [];

        this._initElements();
        if(this.canvas) this._initDefaults();
        console.log('ReportChart canvas:', this.canvas);
    }

    _initElements(){
        this.canvas = document.getElementById(this.canvasId);
        this.startInput = document.getElementById(this.startInputId);
        this.endInput = document.getElementById(this.endInputId);
        this.applyBtn = document.getElementById(this.applyBtnId);
        this.chartType = document.getElementById(this.chartType);
        this.metricInput = document.getElementById(this.metricInputId);
        this.metricTotal = document.getElementById(this.metricTotalId);
        this.average = document.getElementById(this.averageId);

        if(this.canvas) this._ctx = this.canvas.getContext('2d');
    }

    _initDefaults(){
 
        const today = new Date();
        const toISODate = d=>d.toISOString().slice(0,10);
        if(this.endInput) this.endInput.value = toISODate(today);
        if(this.startInput){ const past = new Date(today); past.setDate(past.getDate()-29); this.startInput.value = toISODate(past); }

        if(this.applyBtn) this._addListener(this.applyBtn, 'click', ()=>this.update());
        if(this.chartType) this._addListener(this.chartType, 'change', ()=>this.update());

        // auto-initialize on creation
        this.update();
    }

    _addListener(el, ev, fn){ if(!el) return; el.addEventListener(ev, fn); this._listeners.push({el, ev, fn}); }

    _removeListeners(){ for(const l of this._listeners){ try{ l.el.removeEventListener(l.ev, l.fn); }catch(e){} } this._listeners=[]; }

    async update(){
        console.log('Updating report chart...');
        const urlbase = this.urlbase || '/fin/charts';
        const chartid = this.chartID
        const chartType = this.chartType?.value || 'line';
        const metric = this.metricInput?.value?.trim();
        const start = this.startInput?.value;
        const end = this.endInput?.value;
        // if(!metric){ if(!this.metricInput) return; alert('请输入指标'); return; }
        console.log(`Fetching data for metric=${metric}, start=${start}, end=${end}...`);

        const url = `/api/v1${urlbase}/${encodeURIComponent(chartType)}/${chartid}?start=${encodeURIComponent(start)}&end=${encodeURIComponent(end)}`;
        let data = null;
        try{
            const res = await fetch(url, {
                cache: 'no-store',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });
            if(!res.ok) {
                console.warn('failed to fetch report data, using sample', res.status);
                data = this.sampleData(metric, start, end);
            }else{
                data = await res.json();
            }
        }catch(e){
            console.error('error fetching report data, using sample', e);
            data = this.sampleData(metric, start, end);
        }

        console.log('Report data received:', data);
        console.log(this.sampleData(metric, start, end));

        if(this.metricTotal) this.metricTotal.textContent = data.total ?? '—';

        if(this.average) this.average.textContent = data.avg ?? '—';

        console.log('Rendering chart...');

        const type = this.chartType?.value || 'line';
        if(type === 'line'){
            const series = this.normalizeSeries(data);
            this.renderLine(series, data.title);
        } else {
            const pie = this.normalizePie(data);
            this.renderPie(pie, data.title);
        }

        console.log('Report chart updated.');

        return null;
    }

    normalizeSeries(data){
        if(Array.isArray(data.series)) return data.series.map(s=>({ name: s.name || s.label || 'series', points: (s.points || s.data || []).map(p=>({t: p.t || p.date, v: p.v ?? p.value})) }));
        if(data.points) return [{name: data.name||'series', points: data.points}];
        return [];
    }

    normalizePie(data){
        if(Array.isArray(data.pie)) return data.pie.map(p=>({label:p.label, value:p.value}));
        if(Array.isArray(data.series)) return data.series.map(s=>{ const total = (s.points||s.data||[]).reduce((acc,p)=>acc+(p.v??(p.value||0)),0); return {label: s.name||s.label, value: total}; });
        return [];
    }

    aggregateTotalFromSeries(data){ if(!Array.isArray(data.series)) return '—'; return data.series.flatMap(s => s.points || s.data || []).reduce((acc,p)=>acc + (p.v ?? p.value ?? 0), 0); }

    renderLine(series, title){
        if(!this._ctx) return;
        const datasets = series.map((s,i)=>({ label: s.name || `Series ${i+1}`, data: (s.points||[]).map(p=>({x: this.parseDate(p.t), y: Number(p.v)||0})), borderColor: this.palette(i), backgroundColor: this.hexToRgba(this.palette(i),0.12), tension:0.25, pointRadius:3, fill:true }));
        const config = { type:'line', data:{datasets}, options:{ responsive:true, interaction:{mode:'index', intersect:false}, plugins:{legend:{position:'top'}, title:{display:true, text:`${title}`}}, scales:{ x:{type:'time', time:{unit:'day', tooltipFormat:'yyyy-MM-dd'}, title:{display:true,text:'Date'}}, y:{beginAtZero:true, title:{display:true,text:'Value'}} } } };
        if(this._chart) this._chart.destroy();
        this._chart = new Chart(this._ctx, config);
    }

    renderPie(pieData, title){ if(!this._ctx) return; const labels = pieData.map(p=>p.label); const values = pieData.map(p=>p.value); const colors = labels.map((_,i)=>this.palette(i)); const config = { type:'pie', data:{labels, datasets:[{data: values, backgroundColor: colors}]}, options:{ responsive:true, plugins:{ legend:{position:'right'}, title:{display:true, text:title}, tooltip:{callbacks:{label: ctx => { const val = ctx.parsed; const sum = values.reduce((a,b)=>a+b,0); const pct = sum ? (val/sum*100).toFixed(1)+'%' : '0%'; return `${ctx.label}: ${val} (${pct})`; }}} } } }; if(this._chart) this._chart.destroy(); this._chart = new Chart(this._ctx, config); }

    parseDate(v){ if(!v) return null; if(typeof v === 'number') return new Date(v); if(typeof v === 'string') return new Date(v); return v; }
    palette(i){ const colors=['#60a5fa','#f97316','#a78bfa','#06b6d4','#f59e0b','#ef4444','#34d399','#f472b6']; return colors[i%colors.length]; }
    hexToRgba(hex,a){ const h=hex.replace('#',''); const n=parseInt(h,16); const r=(n>>16)&255,g=(n>>8)&255,b=n&255; return `rgba(${r},${g},${b},${a})`; }

    sampleData(metric, start, end){ const toISODate = d=>d.toISOString().slice(0,10); const s=new Date(start), e=new Date(end); const days=Math.max(1, Math.round((e-s)/(1000*60*60*24))+1); const seriesCount=2; const series=[]; for(let i=0;i<seriesCount;i++){ const points=[]; for(let d=0; d<days; d++){ const dt = new Date(s); dt.setDate(s.getDate()+d); const base = 80*(i+1) + Math.sin(d/5+i)*25; points.push({t: toISODate(dt), v: Math.max(0, Math.round(base + Math.random()*40-20))}); } series.push({name:`Group ${i+1}`, points}); } return {series, total: series.flatMap(s=>s.points).reduce((a,p)=>a+(p.v||0),0), activeUsers: Math.floor(Math.random()*200+50), growth: (Math.random()*10-2).toFixed(1) + '%'}; }

    destroy(){ this._removeListeners(); if(this._chart){ try{ this._chart.destroy(); }catch(e){} this._chart=null; } }

    // helper: if the page still includes old ids, auto-init for backward compatibility
    static autoInitIfPresent(opts){ const canvas = document.getElementById((opts && opts.canvasId) || 'mainChart'); if(!canvas) return null; return new ReportChart(opts); }
}

export default ReportChart;