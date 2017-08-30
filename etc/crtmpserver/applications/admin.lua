application=
{
        name="admin",
        description="Application for administering",
        protocol="dynamiclinklibrary",
        aliases=
        {
                "admin_alias1",
                "admin_alias2",
                "admin_alias3",
        },
        acceptors =
        {
                {
                        ip="0.0.0.0",
                        port=1112,
                        protocol="inboundJsonCli",
                        useLengthPadding=true
                },
        }
        --validateHandshake=true,
        --default=true,
}