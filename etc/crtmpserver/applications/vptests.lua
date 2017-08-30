application=
{
	name="vptests",
	description="Variant protocol tests",
	protocol="dynamiclinklibrary",
	aliases=
	{
		"vptests_alias1",
		"vptests_alias2",
		"vptests_alias3",
	},
	acceptors =
	{
		{
			ip="0.0.0.0",
			port=1111,
			protocol="inboundHttpXmlVariant"
		}
	}
	--validateHandshake=true,
	--default=true,
}