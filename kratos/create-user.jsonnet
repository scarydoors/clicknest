function(ctx) {
    user_id: ctx.identity.id,
    email: ctx.identity.traits.email,
    first_name: ctx.identity.traits.name.first,
    last_name: ctx.identity.traits.name.last,
}
