module.exports = {
    title: "Commercio.network Documentation",
    description: "Documentation for the Commercio.network blockchain.",
    ga: "UA-51029217-2",
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        repo: "commercionetwork/commercionetwork",
        editLinks: true,
        docsDir: "docs",
        docsBranch: "master",
        editLinkText: 'Edit this page on Github',
        lastUpdated: true,
        nav: [
            {text: "Commercio.network", link: "https://commercio.network"},
        ],
        sidebarDepth: 2,
        sidebar: [
            {
                title: "Nodes",
                collapsable: false,
                children: [
                    ["nodes/", "Introduction"],
                    ["nodes/hardware-requirements", "Hardware requirements"],
                    ["nodes/full-node-installation", "Installing a full node"],
                    ["nodes/validator-node-installation", "Becoming a validator"],
                    ["nodes/validator-node-update", "Updating a validator"],
                ]
            },
            {
                title: "Developers",
                collapsable: false,
                children: [
                    ["developers/", "Introduction"],
                    "developers/create-sign-broadcast-tx"
                ]
            },
            {
                title: "Modules",
                collapsable: false,
                children: [
                    "x/government/",
                    "x/docs/",
                    "x/id/"
                ]
            },
            "chain-start/"
        ],
    }
};