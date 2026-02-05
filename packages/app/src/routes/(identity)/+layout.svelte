<script lang="ts">
let { children } = $props()
</script>

<div class="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
	<div class="w-full max-w-sm">
        {@render children()}
	</div>
</div>
