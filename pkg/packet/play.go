package packet

import (
	"io"
)

// IDs mostly match mojang, and are noted when different.
const (
	ClientPlayTeleportConfirmID = iota
	ClientPlayBlockEntityTagQueryID
	ClientPlayChangeDifficultyID
	ClientPlayChatAckID
	ClientPlayChatCommandID
	ClientPlayChatCommandSignedID
	ClientPlayChatID
	ClientPlayChatSessionUpdateID
	ClientPlayChunkBatchReceivedID
	ClientPlayClientStatusID   // Mojang is client command
	ClientPlayClientSettingsID // Mojang is client information
	ClientPlayCommandSuggestionID
	ClientPlayConfigurationAckID
	ClientPlayContainerButtonClickID
	ClientPlayContainerClickID
	ClientPlayContainerCloseID
	ClientPlayContainerSlotStateChangedID
	ClientPlayCookieResponseID
	ClientPlayPluginMessageID // Mojang is custom payload
	ClientPlayDebugSampleSubscriptionID
	ClientPlayEditBookID
	ClientPlayEntityTagQueryID
	ClientPlayInteractID
	ClientPlayJigsawGenerateID
	ClientPlayKeepAliveID
	ClientPlayLockDifficultyID
	ClientPlayMovePlayerPosID
	ClientPlayMovePlayerPosRotID
	ClientPlayMovePlayerRotID
	ClientPlayMovePlayerStatusOnlyID
	ClientPlayMoveVehicleID
	ClientPlayPaddleBoatID
	ClientPlayPickItemID
	ClientPlayPingRequestID
	ClientPlayPlaceRecipeID
	ClientPlayPlayerAbilitiesID
	ClientPlayPlayerActionID
	ClientPlayPlayerCommandID
	ClientPlayPlayerInputID
	ClientPlayPongID
	ClientPlayRecipeBookChangeSettingsID
	ClientPlayRecipeBookSeenRecipeID
	ClientPlayRenameItemID
	ClientPlayResourcePackStatusID // Mojang is just resource pack
	ClientPlaySeenAdvancementsID
	ClientPlaySelectTradeID
	ClientPlaySetBeaconID
	ClientPlaySetCarriedItemID
	ClientPlaySetCommandBlockID
	ClientPlaySetCommandMinecartID
	ClientPlaySetCreativeModeSlotID
	ClientPlaySetJigsawBlockID
	ClientPlaySetStructureBlockID
	ClientPlaySetJigsawID
	ClientPlaySignUpdateID
	ClientPlaySwingID
	ClientPlayTeleportToEntityID
	ClientPlayUseItemOnID
	ClientPlayUseItemID
)

type ClientConfigurationAck struct{}

func (p *ClientConfigurationAck) Direction() Direction { return Serverbound }
func (p *ClientConfigurationAck) ID(state State) int {
	return stateId1(state, Play, ClientPlayConfigurationAckID)
}
func (p *ClientConfigurationAck) Read(_ io.Reader) (err error) {
	return nil
}
func (p *ClientConfigurationAck) Write(_ io.Writer) (err error) {
	return nil
}

// IDs mostly match mojang, and are noted when different.
const (
	ServerPlayBundleDelimiterID = iota
	ServerPlayAddEntityID
	ServerPlayAddExperienceOrbID
	ServerPlayAnimateEntityID // Mojang is just 'animate'
	ServerPlayAwardStatsID
	ServerPlayBlockChangedAckID
	ServerPlayBlockDestructionID
	ServerPlayBlockEntityDataID
	ServerPlayBlockEventID
	ServerPlayBlockUpdateID
	ServerPlayBossBarID // Mojang is boss event
	ServerPlayChangeDifficultyID
	ServerPlayChunkBatchFinishedID
	ServerPlayChunkBatchStartID
	ServerPlayChunkBiomesID
	ServerPlayClearTitleID // Mojang is clear titles
	ServerPlayCommandSuggestionsID
	ServerPlayCommandsID
	ServerPlayContainerCloseID
	ServerPlayContainerSetContentID
	ServerPlayContainerSetDataID
	ServerPlayContainerSetSlotID
	ServerPlayCookieRequestID
	ServerPlayCooldownID
	ServerPlayCustomChatCompletionsID
	ServerPlayPluginMessageID // Mojang is custom payload
	ServerPlayDamageEventID
	ServerPlayDebugSampleID
	ServerPlayDeleteChatID
	ServerPlayDisconnectID
	ServerPlayDisguisedChatID
	ServerPlayEntityEventID
	ServerPlayExplosionID   // Mojang is explode
	ServerPlayForgetChunkID // Mojang is forget level chunk
	ServerPlayGameEventID
	ServerPlayHorseScreenOpenID
	ServerPlayHurtAnimationID
	ServerPlayInitializeBorderID
	ServerPlayKeepAliveID
	ServerPlayChunkDataWithLightID
	ServerPlayWorldEventID    // Mojang is level event
	ServerPlayWorldParticleID // Mojang is level particles
	ServerPlayLightUpdateID
	ServerPlayLoginID
	ServerPlayMapDataID
	ServerPlayMerchantOffersID
	ServerPlayMoveEntityPosID
	ServerPlayMoveEntityPosRotID
	ServerPlayMoveEntityRotID
	ServerPlayMoveVehicleID
	ServerPlayOpenBookID
	ServerPlayOpenScreenID
	ServerPlayOpenSignEditorID
	ServerPlayPingID
	ServerPlayPongResponseID
	ServerPlayPlaceGhostRecipeID
	ServerPlayPlayerAbilitiesID
	ServerPlayPlayerChatID
	ServerPlayPlayerCombatEndID
	ServerPlayPlayerCombatEnterID
	ServerPlayPlayerCombatKillID
	ServerPlayPlayerInfoRemoveID
	ServerPlayPlayerInfoUpdateID
	ServerPlayPlayerLookAtID
	ServerPlayPlayerPositionID
	ServerPlayRecipeID
	ServerPlayRemoveEntitiesID
	ServerPlayRemoveEntityEffectID // Mojang is remove mob effect
	ServerPlayRemoveScoreID
	ServerPlayResourcePackPopID
	ServerPlayResourcePackPushID
	ServerPlayRespawnID
	ServerPlayRotateHeadID
	ServerPlaySectionBlocksUpdateID
	ServerPlaySelectAdvancementTabID
	ServerPlayServerDataID
	ServerPlaySetActionBarTextID
	ServerPlaySetWorldCenterID
	ServerPlaySetWorldLerpSizeID
	ServerPlaySetWorldSizeID
	ServerPlaySetWorldWarningDelayID
	ServerPlaySetWorldWarningReachID
	ServerPlaySetCameraID
	ServerPlaySetCarriedItemChangeID
	ServerPlaySetChunkCacheCenterID
	ServerPlaySetChunkCacheRadiusID
	ServerPlaySetDefaultSpawnPositionID
	ServerPlaySetDisplayObjectiveID
	ServerPlaySetEntityDataID
	ServerPlaySetEntityLinkID
	ServerPlaySetEntityVelocityID // Mojang is set entity motion
	ServerPlaySetEquipmentID
	ServerPlaySetExperienceID
	ServerPlaySetHealthID
	ServerPlaySetObjectiveID
	ServerPlaySetPassengersID
	ServerPlaySetPlayerTeamID
	ServerPlaySetScoreID
	ServerPlaySetSimulationDistanceID
	ServerPlaySetSubtitleTextID
	ServerPlaySetTimeID
	ServerPlaySetTitleTextID
	ServerPlaySetTitleTimeID // Mojang is set titles animation
	ServerPlaySoundEntityID
	ServerPlaySoundID
	ServerPlayStartConfigurationID
	ServerPlayStopSoundID
	ServerPlayStoreCookieID
	ServerPlaySystemChatID
	ServerPlayTabListID
	ServerPlayTagQueryID
	ServerPlayTakeItemEntityID
	ServerPlayTeleportEntityID
	ServerPlayTickingStateID
	ServerPlayTickingStepID
	ServerPlayTransferID
	ServerPlayUpdateAdvancementsID
	ServerPlayUpdateEntityAttributesID // Mojang is update attributes
	ServerPlayUpdateEntityEffectID     // Mojang is update entity effect
	ServerPlayUpdateRecipesID
	ServerPlayUpdateTagsID
	ServerPlayProjectilePowerID
)

type ServerStartConfiguration struct {
}

func (p *ServerStartConfiguration) Direction() Direction { return Clientbound }
func (p *ServerStartConfiguration) ID(state State) int {
	return stateId1(state, Play, ServerPlayStartConfigurationID)
}
func (p *ServerStartConfiguration) Read(_ io.Reader) (err error) {
	return nil
}
func (p *ServerStartConfiguration) Write(_ io.Writer) (err error) {
	return nil
}

var (
	_ Packet = (*ClientConfigurationAck)(nil)

	_ Packet = (*ServerStartConfiguration)(nil)
)
