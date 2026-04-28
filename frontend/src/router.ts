// Generouted, changes to this file will be overridden
/* eslint-disable */

import { components, hooks, utils } from '@generouted/react-router/client'

export type Path =
  | `/`
  | `/login`
  | `/login/callback`
  | `/my`
  | `/my/tsumikis`
  | `/my/works`
  | `/tsumikis`
  | `/tsumikis/:tsumikiId`
  | `/upload`
  | `/users/:userId`
  | `/users/:userId/tsumikis`
  | `/users/:userId/works`
  | `/works`
  | `/works/:workId`

export type Params = {
  '/tsumikis/:tsumikiId': { tsumikiId: string }
  '/users/:userId': { userId: string }
  '/users/:userId/tsumikis': { userId: string }
  '/users/:userId/works': { userId: string }
  '/works/:workId': { workId: string }
}

export type ModalPath = never

export const { Link, Navigate } = components<Path, Params>()
export const { useModals, useNavigate, useParams } = hooks<Path, Params, ModalPath>()
export const { redirect } = utils<Path, Params>()
