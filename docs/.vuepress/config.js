module.exports = {
    title: "Commercio.network Documentation",
    description: "Documentation for the Commercio.network blockchain.",
    ga: "UA-51029217-2",
    head: [
        ['link', {rel: 'icon', href: '/icon.png'}]
    ],
    markdown: {
        lineNumbers: true,
	    extendMarkdown: md => {
		    md.use(require("markdown-it-footnote"));
  	    }
    },
    plugins: [
        'latex',
        [
          "sitemap",
          {
            hostname: "https://docs.commercio.network"
          }
        ]
    ],
    themeConfig: {
        repo: "commercionetwork/commercionetwork",
        editLinks: true,
        docsDir: "docs",
        docsBranch: "master",
        editLinkText: 'Edit this page on Github',
        lastUpdated: true,
        //logo: '/.vuepress/icon.png',
        nav: [
            {text: "Commercio.network", link: "https://commercio.network"},
            {text: "ver 2.2.0", link: "/" },
            {text: "ver 2.1.2", link: "/docs2.1.2/" },
        ],
        sidebarDepth: 3,
        sidebar: [
            {
                title: "Running Nodes",
                collapsable: false,
                children: [
                    ["nodes/", "Commercio.network overview"],
                    ["nodes/hardware-requirements", "Hardware requirements"],
                    ["nodes/full-node-installation", "Installing a full node"],
                    ["nodes/validator-node-installation", "Becoming a validator"],
                    ["nodes/validator-node-handling", "Handling a validator"],
                    ["nodes/validator-node-update", "Updating a validator"]
                ]
            },
            {
                title: "API Developers",
                collapsable: false,
                children: [
                    ["app_developers/commercioapi-introduction", "CommercioAPI Introduction"],
                    ["app_developers/commercioapi-authentication", "CommercioAPI Authentication"],
                    ["app_developers/commercioapi-sharedoc", "CommercioAPI ShareDoc"],
                ]
            },
            {
                title: "SDK Developers",
                collapsable: false,
                children: [
                    ["developers/", "Introduction"],
                    "developers/create-sign-broadcast-tx",
                    "developers/message-types",
                    "developers/listing-transactions"
                ]
            },
            {
                title: "Custom Commercio.network Modules",
                collapsable: false,
                children: [
                    "x/bank/",
                    "x/government/",
                    "x/id/",
                    "x/documents/",
                    "x/commerciomint/",
                    "x/commerciokyc/",
                    "x/vbr/"
                ]
            },
            {
                title: "ver 2.1.2",
                collapsable: true,
                children: [
                    ["docs2.1.2/", "ver 2.1.2"],
                    {
                        title: "Nodes",
                        collapsable: true,
                        children: [
                            ["docs2.1.2/nodes/", "Introduction"],
                            ["docs2.1.2/nodes/hardware-requirements", "Hardware requirements"],
                            ["docs2.1.2/nodes/full-node-installation", "Installing a full node"],
                            ["docs2.1.2/nodes/validator-node-installation", "Becoming a validator"],
                            ["docs2.1.2/nodes/validator-node-handling", "Handling a validator"],
         //                   ["docs2.1.2/nodes/validator-node-installation-mainnet", "Becoming a validor in the Mainnet"],
                            ["docs2.1.2/nodes/validator-node-update", "Updating a validator"],
                        ]
                    },
                    {
                        title: "App Developers",
                        collapsable: true,
                        children: [
                            ["docs2.1.2/app_developers/", "Introduction"]
                        ]
                    },
                    {
                        title: "SDK Developers",
                        collapsable: true,
                        children: [
                            ["docs2.1.2/developers/", "Introduction"],
                            "docs2.1.2/developers/create-sign-broadcast-tx",
                            "docs2.1.2/developers/message-types",
                            "docs2.1.2/developers/listing-transactions"
                        ]
                    },


                    {
                        title: "Modules",
                        collapsable: true,
                        children: [
                            "docs2.1.2/x/bank/",
                            "docs2.1.2/x/government/",
                            "docs2.1.2/x/id/",
                            "docs2.1.2/x/docs/",
                            "docs2.1.2/x/pricefeed/",
                            "docs2.1.2/x/commerciomint/",
                            "docs2.1.2/x/memberships/",
                            "docs2.1.2/x/vbr/",
                            "docs2.1.2/x/creditrisk/"
                        ]
                    }
                ]
            }

        ],
    }
};
