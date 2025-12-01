#!/bin/bash

echo "🧹 Cleaning Nuxt cache..."
rm -rf .nuxt .output node_modules/.cache

echo "📦 Ensuring public directory exists..."
mkdir -p public
mkdir -p public/assets
mkdir -p public/images

# Create placeholder if public directory is empty
if [ -z "$(ls -A public)" ]; then
  echo "📝 Creating placeholder file in public directory..."
  echo "/* MSX Player Static Assets */" > public/placeholder.txt
fi

echo "🚀 Starting Nuxt development server with Deno 3.0..."
deno run -A npm:nuxt dev --host localhost --port 3000