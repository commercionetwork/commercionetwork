module.exports = {
    title: "Commercio.network SDK Documentation",
    description: "Documentation for the Commercio.network blockchain.",
    ga: "UA-51029217-2",
    // dest: "./dist/docs",
    // base: "/docs/",
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
            { text: "Back to Commercio.network", link: "https://commercio.network" },
        ],
        sidebar: {
            "/" : [
                "/installation",
                "/join-testnet",
                "/join-mainnet",
                "/validator-setup",
                "/validator-hardware",
            ]
        }
    }
};