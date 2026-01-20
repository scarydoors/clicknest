<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
    import * as Sidebar from '$lib/components/ui/sidebar/index.ts';
    import AppSidebar from '$lib/components/app-sidebar.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}
