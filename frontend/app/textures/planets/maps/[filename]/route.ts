import { promises as fs } from 'node:fs'
import path from 'node:path'

const DEFAULT_PROJECT_ROOT = '/home/weaver/learnweaver'

function getPlanetTextureMapDir(): string {
  const explicitDir = process.env.PLANET_TEXTURE_MAP_DIR?.trim()
  if (explicitDir) {
    return explicitDir
  }

  const projectRoot =
    process.env.LEARNWEAVER_PROJECT_ROOT?.trim() ||
    process.env.PROJECT_ROOT?.trim() ||
    DEFAULT_PROJECT_ROOT

  return path.join(projectRoot, 'frontend', 'public', 'textures', 'planets', 'maps')
}

export const runtime = 'nodejs'

function isSafePlanetTextureMapFilename(filename: string): boolean {
  return /^[0-9a-f-]+\.webp$/i.test(filename) || filename === 'learnweaver_planet_atlas_basic_seamfixed_2048x768.webp'
}

export async function GET(
  _request: Request,
  context: { params: Promise<{ filename: string }> },
) {
  const { filename } = await context.params

  if (!isSafePlanetTextureMapFilename(filename)) {
    return new Response('Not found', { status: 404 })
  }

  const filePath = path.join(getPlanetTextureMapDir(), filename)

  try {
    const buffer = await fs.readFile(filePath)
    return new Response(buffer, {
      status: 200,
      headers: {
        'Content-Type': 'image/webp',
        'Cache-Control': 'public, max-age=60, stale-while-revalidate=300',
      },
    })
  } catch (error) {
    if ((error as NodeJS.ErrnoException).code === 'ENOENT') {
      return new Response('Not found', { status: 404 })
    }
    return new Response('Failed to read asset', { status: 500 })
  }
}
