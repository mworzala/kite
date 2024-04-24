package packet

import (
	"io"
)

const (
	ClientPlayTeleportConfirmID = iota
	ClientPlayQueryBlockNbtID
	ClientPlayDifficultyID
	ClientPlayChatAckID
	ClientPlayCommandChatID
	ClientPlaySignedCommandChatID
	ClientPlayChatMessageID
	ClientPlayChatSessionUpdateID
	ClientPlayChunkBatchReceivedID
	ClientPlayStatusID
	ClientPlaySettingsID
	ClientPlayTabCompleteID
	ClientPlayConfigurationAckID
)

type ClientPlayConfigurationAck struct{}

func (p *ClientPlayConfigurationAck) Direction() Direction { return Serverbound }
func (p *ClientPlayConfigurationAck) ID(state State) int {
	return stateId1(state, Play, ClientPlayConfigurationAckID)
}

func (p *ClientPlayConfigurationAck) Read(r io.Reader) (err error) {
	return nil
}

func (p *ClientPlayConfigurationAck) Write(w io.Writer) (err error) {
	return nil
}

const (
	ServerPlayBundleDelimiterID = iota
	ServerPlaySpawnEntityID
	ServerPlaySpawnExperienceOrbID
	ServerPlayEntityAnimationID
	ServerPlayStatisticsID
	ServerPlayAcknowledgeBlockChangeID
	ServerPlayBlockBreakAnimationID
	ServerPlayBlockEntityDataID
	ServerPlayBlockActionID
	ServerPlayBlockChangeID
	ServerPlayBossBarID
	ServerPlayDifficultyID
	ServerPlayChunkBatchFinishedID
	ServerPlayChunkBatchStartID
	ServerPlayChunkBiomesID
	ServerPlayClearTitleID
	ServerPlayTabCompleteID
	ServerPlayDeclareCommandsID
	ServerPlayCloseWindowID
	ServerPlayWindowItemsID
	ServerPlayWindowPropertyID
	ServerPlaySetSlotID
	ServerPlayCookieRequestID
	ServerPlaySetCooldownID
	ServerPlayCustomChatCompletionsID
	ServerPlayPluginMessageID
	ServerPlayDamageEventID
	ServerPlayDebugSampleID
	ServerPlayDeleteChatMessageID
	ServerPlayDisconnectID
	ServerPlayDisguisedChatID
	ServerPlayEntityStatusID
	ServerPlayExplosionID
	ServerPlayUnloadChunkID
	ServerPlayChangeGameStateID
	ServerPlayOpenHorseWindowID
	ServerPlayHitAnimationID
	ServerPlayInitWorldBorderID
	ServerPlayKeepAliveID
	ServerPlayChunkDataID
	ServerPlayEffectID
	ServerPlayParticleID
	ServerPlayUpdateLightID
	ServerPlayJoinGameID
	ServerPlayMapDataID
	ServerPlayTradeListID
	ServerPlayEntityPositionID
	ServerPlayEntityPositionAndRotationID
	ServerPlayEntityRotationID
	ServerPlayVehicleMoveID
	ServerPlayOpenBookID
	ServerPlayOpenWindowID
	ServerPlayOpenSignEditorID
	ServerPlayPingID
	ServerPlayPingResponseID
	ServerPlayCraftRecipeResponseID
	ServerPlayPlayerAbilitiesID
	ServerPlayPlayerChatID
	ServerPlayEndCombatEventID
	ServerPlayEnterCombatEventID
	ServerPlayPlayerInfoRemoveID
	ServerPlayPlayerInfoUpdateID
	ServerPlayFacePlayerID
	ServerPlayPlayerPositionAndLookID
	ServerPlayUnlockRecipesID
	ServerPlayDestroyEntitiesID
	ServerPlayRemoveEntityEffectID
	ServerPlayResetScoreID
	ServerPlayResourcePackPushID
	ServerPlayResourcePackPopID
	ServerPlayRespawnID
	ServerPlayEntityHeadLookID
	ServerPlayMultiBlockChangeID
	ServerPlaySelectAdvancementTabID
	ServerPlayServerDataID
	ServerPlayActionBarID
	ServerPlayWorldBorderCenterID
	ServerPlayWorldBorderLerpSizeID
	ServerPlayWorldBorderSizeID
	ServerPlayWorldBorderWarningDelayID
	ServerPlayWorldBorderWarningReachID
	ServerPlayCameraID
	ServerPlayHeldItemChangeID
	ServerPlayUpdateViewPositionID
	ServerPlayUpdateViewDistanceID
	ServerPlaySpawnPositionID
	ServerPlayDisplayScoreboardID
	ServerPlayEntityMetadataID
	ServerPlayAttachEntityID
	ServerPlayEntityVelocityID
	ServerPlayEntityEquipmentID
	ServerPlaySetExperienceID
	ServerPlayUpdateHealthID
	ServerPlayScoreboardObjectiveID
	ServerPlaySetPassengersID
	ServerPlayTeamsID
	ServerPlayUpdateScoreID
	ServerPlaySetSimulationDistanceID
	ServerPlaySetTitleSubtitleID
	ServerPlayTimeUpdateID
	ServerPlaySetTitleTextID
	ServerPlaySetTitleTimeID
	ServerPlayEntitySoundEffectID
	ServerPlaySoundEffectID
	MISSING_ONE_HERE_IDK_WHAT_IT_IS
	ServerPlayStartConfigurationID
	ServerPlayStopSoundID
	ServerPlayCookieStoreID
	ServerPlaySystemChatID
	ServerPlayPlayerListHeaderAndFooterID
	ServerPlayNbtQueryResponseID
	ServerPlayCollectItemID
	ServerPlayEntityTeleportID
	ServerPlayTickStateID
	ServerPlayTickStepID
	ServerPlayTransferID
	ServerPlayAdvancementsID
	ServerPlayEntityAttributesID
	ServerPlayEntityEffectID
	ServerPlayDeclareRecipesID
	ServerPlayUpdateTagsID
	ServerPlayProjectilePowerID
)

type ServerPlayStartConfiguration struct {
}

func (p *ServerPlayStartConfiguration) Direction() Direction { return Serverbound }
func (p *ServerPlayStartConfiguration) ID(state State) int {
	return stateId1(state, Play, ServerPlayStartConfigurationID)
}

func (p *ServerPlayStartConfiguration) Read(r io.Reader) (err error) {
	return nil
}

func (p *ServerPlayStartConfiguration) Write(w io.Writer) (err error) {
	return nil
}

var (
	_ Packet = (*ClientPlayConfigurationAck)(nil)

	_ Packet = (*ServerPlayStartConfiguration)(nil)
)
