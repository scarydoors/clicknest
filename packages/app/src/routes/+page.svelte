<script lang="ts">
import * as Card from "$lib/components/ui/card/index"
import * as Chart from "$lib/components/ui/chart/index";
import { AreaChart } from "layerchart";
import { scaleUtc } from "d3-scale";
	import { curveLinear } from "d3-shape";

const chartData = [
    { date: new Date("2024-01-01"), value: 186 },
    { date: new Date("2024-02-01"), value: 305 },
    { date: new Date("2024-03-01"), value: 237 },
    { date: new Date("2024-04-01"), value: 73 },
    { date: new Date("2024-05-01"), value: 209 },
    { date: new Date("2024-06-01"), value: 214 },
];

const chartConfig = {
    value: {label:"Pageviews", color: "var(--chart-1)"}
} satisfies Chart.ChartConfig;
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<Chart.Container config={chartConfig}>
    <AreaChart
        x="date"
        series={[
            {
                key: "value",
                label: "Pageviews",
                color: chartConfig.value.color
            }
        ]}
        data={chartData}
        xScale={scaleUtc()}
        props={{
            area: {
                curve: curveLinear,
                "fill-opacity": 0.4,
                line: { class: "stroke-1" },
                motion: "tween",
            },
            xAxis: {
                format: (v: Date) => v.toLocaleDateString(undefined, { month: "short" })
            },
            yAxis: { format: () => "" },
        }}
    >
        {#snippet tooltip()}
            <Chart.Tooltip hideLabel />
        {/snippet}
    </AreaChart>
</Chart.Container>
