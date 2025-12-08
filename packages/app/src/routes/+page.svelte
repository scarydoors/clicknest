<script lang="ts">
import * as Card from "$lib/components/ui/card/index"
import * as Chart from "$lib/components/ui/chart/index";
import { AreaChart } from "layerchart";
import { scaleUtc } from "d3-scale";
import { curveLinear } from "d3-shape";
import { onMount } from "svelte";

const chartConfig = {
    value: {label:"Pageviews", color: "var(--chart-1)"}
} satisfies Chart.ChartConfig;

type Timeseries = { timestamp: Date, value: number }[]

let loading = $state<boolean>(true);
let data = $state<Timeseries | undefined>();

onMount(() => {
    let yesterday = new Date();
    yesterday.setUTCDate(6);
    let tomorrow = new Date();
    tomorrow.setUTCDate(9)
    fetch("http://localhost:6969/api/timeseries?" + new URLSearchParams({
        "domain": "stupidwebsite.com",
        "interval": "1m",
        "start_date": yesterday.toISOString(),
        "end_date": tomorrow.toISOString(),

    }).toString()).then(async (resp) => {
        const json = await resp.json();
        data = json.map((t: {timestamp: string, value: number}) => ({ ...t, timestamp: new Date(t.timestamp) }))
        loading = false;
    })
})
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>

<Card.Root>
    <Chart.Container config={chartConfig}>
        <AreaChart
            x="timestamp"
            series={[
                {
                    key: "value",
                    label: "Pageviews",
                    color: chartConfig.value.color
                }
            ]}
            data={data}
            xScale={scaleUtc()}
            props={{
                area: {
                    curve: curveLinear,
                    "fill-opacity": 0.4,
                    line: { class: "stroke-1" },
                },
                xAxis: {
                    format: (v: Date) => v.toLocaleDateString(undefined, { month: "short", day: "numeric" }),
                    ticks: data?.length ?? 0
                },
                yAxis: { format: () => "" },
            }}
        >
            {#snippet tooltip()}
                <Chart.Tooltip hideLabel />
            {/snippet}
        </AreaChart>
    </Chart.Container>
</Card.Root>
