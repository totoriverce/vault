import { productName, productSlug } from 'data/metadata'
import DocsPage from '@hashicorp/react-docs-page'
import {
  generateStaticPaths,
  generateStaticProps,
} from '@hashicorp/react-docs-page/server'

const NAV_DATA_FILE = 'data/guides-nav-data.json'
const CONTENT_DIR = 'content/guides'
const basePath = 'guides'

export default function GuidesLayout(props) {
  return (
    <DocsPage
      product={{ name: productName, slug: productSlug }}
      baseRoute={basePath}
      staticProps={props}
    />
  )
}

export async function getStaticPaths() {
  return {
    fallback: false,
    paths: await generateStaticPaths({
      navDataFile: NAV_DATA_FILE,
      localContentDir: CONTENT_DIR,
    }),
  }
}

export async function getStaticProps({ params }) {
  return {
    props: await generateStaticProps({
      mainBranch: 'master',
      navDataFile: NAV_DATA_FILE,
      localContentDir: CONTENT_DIR,
      product: { name: productName, slug: productSlug },
      params,
    }),
  }
}
