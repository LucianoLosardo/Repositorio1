#!/bin/bash

# Script para iniciar los servicios de Docker Compose 

# Opciones para docker-compose:
# -d: Modo "detached" (ejecuta en segundo plano)
# --build: Reconstruye las imágenes si hay cambios en el Dockerfile o contexto.
COMPOSE_OPTIONS="up --build -d"

#--generación de código templ--#
templ generate

# --- Detección del comando Compose (V2 vs V1) ---
CMD=""
if docker compose version &>/dev/null; then
    # El nuevo 'docker compose' (V2) existe
    CMD="docker compose"
elif docker-compose version &>/dev/null; then
    # El viejo 'docker-compose' (V1) existe
    CMD="docker-compose"
else
    echo "Error: No se encontró ni 'docker compose' (v2) ni 'docker-compose' (v1)."
    echo "Por favor, asegúrate de que Docker Compose esté instalado."
    exit 1
fi

echo "Usando comando: $CMD"


# --- Verificación de permisos (el problema de tu log) ---
SUDO_PREFIX=""

# Comprueba si el usuario NO es root (EUID != 0) Y NO está en el grupo 'docker'
if [ "$EUID" -ne 0 ] && ! groups $USER | grep -q '\bdocker\b'; then
    echo "Aviso: Tu usuario no pertenece al grupo 'docker'. Se intentará usar 'sudo'."
    echo "       (Para evitar esto, ejecuta: sudo usermod -aG docker \$USER y reinicia tu sesión)"
    SUDO_PREFIX="sudo"
fi

# --- Ejecución ---
echo "--------------------------------------------------"
echo "Levantando servicios con: $SUDO_PREFIX $CMD $COMPOSE_OPTIONS"
echo "--------------------------------------------------"

# Ejecuta el comando final
if [ -n "$SUDO_PREFIX" ]; then
    # Si SUDO_PREFIX no está vacío, úsalo
    $SUDO_PREFIX $CMD $COMPOSE_OPTIONS
else
    # Si está vacío, ejecuta el comando directamente
    $CMD $COMPOSE_OPTIONS
fi

# Asegura que el script falle si un comando falla
set -e


